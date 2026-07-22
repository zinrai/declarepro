package main

import (
	"flag"
	"fmt"
	"os"
)

// declarepro drives reprepro declaratively from a YAML config.
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `declarepro - declarative reprepro driver

Usage:
  declarepro [-config declarepro.yaml] [-dry-run]

Flags:
  -config <path>   Path to the YAML config (default "declarepro.yaml")
  -dry-run         Show what would change without modifying anything
  -version         Print version and exit
`)
}

func run() error {
	flag.Usage = usage
	configPath := flag.String("config", "declarepro.yaml", "path to the YAML config")
	dryRun := flag.Bool("dry-run", false, "show what would change without modifying anything")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		printVersion()
		return nil
	}

	c, err := loadConfig(*configPath)
	if err != nil {
		return err
	}
	return converge(c, *dryRun)
}
