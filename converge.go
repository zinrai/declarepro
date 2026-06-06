package main

import "fmt"

func converge(c *Config, dryRun bool) error {
	r := &Runner{DryRun: dryRun}

	if dryRun {
		fmt.Println("== dry-run ==")
	} else {
		fmt.Println("== apply ==")
	}

	if err := convergeConf(c, dryRun); err != nil {
		return err
	}
	return convergePackages(c, r)
}

func convergeConf(c *Config, dryRun bool) error {
	fmt.Println("\n-- config --")

	desired := renderConf(c)
	current, err := readConf(c.BaseDir)
	if err != nil {
		return err
	}

	if current == desired {
		fmt.Println("  conf/distributions: up to date")
		return nil
	}

	if current == "" {
		fmt.Println("  conf/distributions: will be created")
	} else {
		fmt.Println("  conf/distributions: will be updated")
	}
	fmt.Print(confDiff(current, desired))

	if dryRun {
		return nil
	}

	fmt.Println("  writing conf/distributions")
	return writeConf(c.BaseDir, desired)
}

func convergePackages(c *Config, r *Runner) error {
	fmt.Println("\n-- packages --")

	for _, dist := range c.distNames() {
		if err := convergeDistPackages(c, dist, r); err != nil {
			return err
		}
	}
	return nil
}

func convergeDistPackages(c *Config, dist string, r *Runner) error {
	fmt.Printf("\n[%s]\n", dist)

	pkgs := c.Distributions[dist].Packages
	if len(pkgs) == 0 {
		fmt.Println("  no packages declared")
		return nil
	}

	for _, pkg := range pkgs {
		if err := r.repreproIncludeDeb(c.BaseDir, dist, pkg.File); err != nil {
			return err
		}
	}
	return nil
}
