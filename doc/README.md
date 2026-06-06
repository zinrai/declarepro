# declarepro deployment guide

This guide walks through a realistic deployment shape: one reprepro
repository per OS family, managed by one declarepro config each. The shape
maps directly onto the apt-line layout below. The accompanying
`debian.yaml` and `ubuntu.yaml` are the configs the procedures here refer
to.

## Target apt-line

```
deb http://repo.example.internal/debian trixie main
deb http://repo.example.internal/ubuntu noble  main
```

The `/debian` and `/ubuntu` URL paths come from two separate `base_dir`s
served by the same web server. They are independent reprepro repositories;
nothing in declarepro ties them together.

## Files

- `debian.yaml` — Debian side. `base_dir: /srv/repo/debian`, one entry per
  Debian codename (`bookworm`, `trixie`, ...).
- `ubuntu.yaml` — Ubuntu side. `base_dir: /srv/repo/ubuntu`, one entry per
  Ubuntu codename (`jammy`, `noble`, ...).

Each config carries the per-distribution `conf/distributions` stanza
verbatim. Reading the YAML tells you exactly what reprepro will see.

## Running

Apply each config independently:

    declarepro -config debian.yaml
    declarepro -config ubuntu.yaml

There is no combined invocation. Wiring the two together is the operator's
job: an Ansible playbook with two tasks, a Makefile with two targets, a
CI pipeline with two steps. declarepro stays narrowly scoped to one
`base_dir` per run.

Order between the two does not matter: the configs touch disjoint
directories and disjoint reprepro databases.

## Adding a distribution

Add a new entry under `distributions:` and re-apply. reprepro will create
the new distribution on first `includedeb`. Example: adding `bookworm` to
`debian.yaml`:

```yaml
distributions:
  bookworm:                # new
    config: |
      Codename: bookworm
      Architectures: amd64
      Components: main
      SignWith: ABCD1234
    packages:
      - file: /srv/debs/debian/bookworm/mytool_0.9_amd64.deb
  trixie:
    ...
```

## Removing a distribution

Delete the entry from the YAML and re-apply. `conf/distributions` will be
rewritten without it.

The reprepro database keeps the dropped distribution's entries. **This
prevents reprepro from operating on the base directory at all** — including
`includedeb` against the distributions that remain — until the orphan
entries are removed. The next apply will fail with output like:

    Error: packages database contains unused 'bookworm|main|amd64' database.

Cleanup procedure:

    reprepro -b /srv/repo/debian --delete clearvanished

This drops the orphaned database entries. The on-disk `dists/<codename>/`
directory is not removed by `clearvanished`; delete it by hand if you want
the artifacts gone:

    rm -rf /srv/repo/debian/dists/bookworm

After cleanup, re-apply declarepro normally.

declarepro never issues removes itself; the operator decides when to clean
up dropped distributions.

## Version changes

For a package already registered in a distribution:

- **Upgrade** (newer version): reprepro replaces the registered version and
  removes the old `.deb` from the pool automatically. declarepro just calls
  `includedeb` on the new file; no remove is issued.

- **Downgrade** (older version): reprepro silently skips the inclusion as a
  no-op and exits 0. The output reads:

      Skipping inclusion of 'gh' '2.40.0' in 'trixie|main|amd64', as it has already '2.50.0'.

  declarepro returns success. **The downgrade intent is not realized** — the
  database stays at the newer version. To intentionally downgrade, remove
  the existing entry first:

      reprepro -b /srv/repo/debian remove trixie gh

  then re-apply declarepro with the older version in the config.

## Web server

A typical nginx layout:

```nginx
location /debian/ {
    alias /srv/repo/debian/;
    autoindex on;
}
location /ubuntu/ {
    alias /srv/repo/ubuntu/;
    autoindex on;
}
```

## GPG signing

`SignWith: <keyid>` in each stanza tells reprepro to sign `Release` /
`InRelease` with the given GPG key. The key must be importable by the user
running declarepro (i.e. present in that user's GPG keyring). declarepro
does not manage GPG keys; provision them out of band (Ansible, manual
import).

`SignWith` is read by reprepro at export time. **Adding `SignWith` to an
already-populated distribution does not by itself trigger re-signing of the
existing `Release` file**, because reprepro only re-exports when content
changes. Either include a new or updated package to trigger a re-export, or
force one explicitly:

    reprepro -b /srv/repo/debian export trixie

To run without signing, remove the `SignWith` line. reprepro will then
produce unsigned `Release` files. Apt clients will need `[trusted=yes]` in
their sources.list entry to accept an unsigned repository.

## .deb file layout

The examples assume one directory per OS-codename under `/srv/debs/`.
This is one convention; declarepro takes any absolute path and does not
care how files are organized.
