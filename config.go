package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	BaseDir       string                  `yaml:"base_dir"`
	Distributions map[string]Distribution `yaml:"distributions"`
}

type Distribution struct {
	// Config is the verbatim text of this distribution's stanza in
	// conf/distributions. It is written as-is after newline normalization.
	Config string `yaml:"config"`

	Packages []Package `yaml:"packages"`
}

type Package struct {
	File string `yaml:"file"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Config) validate() error {
	if strings.TrimSpace(c.BaseDir) == "" {
		return fmt.Errorf("base_dir is required")
	}
	if len(c.Distributions) == 0 {
		return fmt.Errorf("at least one distribution is required")
	}

	for name, dist := range c.Distributions {
		if err := validateDistribution(name, dist); err != nil {
			return err
		}
	}

	return nil
}

func validateDistribution(name string, dist Distribution) error {
	if strings.TrimSpace(dist.Config) == "" {
		return fmt.Errorf("distribution %q: config is required", name)
	}
	for i, pkg := range dist.Packages {
		if strings.TrimSpace(pkg.File) == "" {
			return fmt.Errorf("distribution %q: packages[%d]: file is required", name, i)
		}
	}
	return nil
}

// distNames returns distribution names in sorted order so output is
// deterministic regardless of map iteration order.
func (c *Config) distNames() []string {
	names := make([]string, 0, len(c.Distributions))
	for name := range c.Distributions {
		names = append(names, name)
	}
	sortStrings(names)
	return names
}
