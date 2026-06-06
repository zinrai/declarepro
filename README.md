# declarepro

Drives [reprepro](https://tracker.debian.org/pkg/reprepro) from a declarative YAML config.

## What it does

declarepro reads a YAML config that declares, for each distribution, the
verbatim `conf/distributions` stanza and the set of `.deb` files that must be
present. It then generates and places `conf/distributions` and runs
`reprepro includedeb` once per declared package. It does not operate
reprepro's database directly; the database is left to reprepro. Every
reprepro command and its output is printed so you can see exactly what ran.

declarepro is a thin translator between the config and reprepro. It does not
pre-classify includedeb as no-op, new install, or version bump; reprepro
makes that decision and reports it. Version bumps are handled by reprepro
itself: including a newer `.deb` into the same distribution replaces the
older version. Packages dropped from the config are not removed.

## Usage

Write a config (see `doc/` for a multi-OS layout):

    base_dir: /srv/repo

    distributions:
      noble:
        config: |
          Codename: noble
          Architectures: amd64
          Components: main
          SignWith: ABCD1234
        packages:
          - file: /srv/deb-archive/mytool_2.1.0_amd64.deb

Preview what would change without modifying anything:

    declarepro -config declarepro.yaml -dry-run

Apply:

    declarepro -config declarepro.yaml

Example apply output:

    == apply ==

    -- config --
      conf/distributions: up to date

    -- packages --

    [noble]
      $ reprepro -b /srv/repo includedeb noble /srv/deb-archive/mytool_2.1.0_amd64.deb

In `-dry-run`, commands that would run are shown with a `+` prefix; without
it, commands that ran are shown with a `$` prefix and their output is
indented below with `|`.

## Behavior notes

`conf/distributions` is managed only through declarepro. Editing it by hand
and then re-running declarepro will overwrite the hand edit.

declarepro does not guard against `conf/distributions` changes that reprepro
cannot handle (for example, removing an architecture while packages for it
remain in the database). The file is written as declared; reprepro will
report any resulting inconsistency on its next invocation, and the user
reconciles by adjusting the config or operating reprepro directly.

## Requirements

`reprepro` must be available on PATH.

## License

This project is licensed under the [MIT License](./LICENSE).
