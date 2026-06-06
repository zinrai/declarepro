package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func confPath(baseDir string) string {
	return filepath.Join(baseDir, "conf", "distributions")
}

// renderConf assembles conf/distributions from the per-distribution stanzas.
// Stanzas are separated by a blank line (reprepro's format) and each is
// trimmed so byte comparison against the on-disk file is stable.
func renderConf(c *Config) string {
	var b strings.Builder
	for i, name := range c.distNames() {
		stanza := normalizeStanza(c.Distributions[name].Config)
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(stanza)
		b.WriteString("\n")
	}
	return b.String()
}

func normalizeStanza(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
	}
	for len(lines) > 0 && lines[0] == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

// readConf returns the on-disk conf/distributions content, or empty string
// if the file does not exist yet.
func readConf(baseDir string) (string, error) {
	data, err := os.ReadFile(confPath(baseDir))
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("reading conf/distributions: %w", err)
	}
	return string(data), nil
}

func writeConf(baseDir, content string) error {
	dir := filepath.Join(baseDir, "conf")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating conf directory: %w", err)
	}
	if err := os.WriteFile(confPath(baseDir), []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing conf/distributions: %w", err)
	}
	return nil
}

// confDiff is a deliberately small line diff: conf/distributions is short
// and the goal is to let a human eyeball the change, not to be a full
// diff engine.
func confDiff(current, desired string) string {
	cur := strings.Split(strings.TrimRight(current, "\n"), "\n")
	des := strings.Split(strings.TrimRight(desired, "\n"), "\n")

	curSet := map[string]int{}
	for _, ln := range cur {
		curSet[ln]++
	}
	desSet := map[string]int{}
	for _, ln := range des {
		desSet[ln]++
	}

	var b strings.Builder
	for _, ln := range cur {
		if desSet[ln] == 0 {
			fmt.Fprintf(&b, "    - %s\n", ln)
		}
	}
	for _, ln := range des {
		if curSet[ln] == 0 {
			fmt.Fprintf(&b, "    + %s\n", ln)
		}
	}
	return b.String()
}
