package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Runner struct {
	DryRun bool
}

type cmdResult struct {
	Stdout string
	Stderr string
}

func (r *Runner) run(name string, args ...string) (cmdResult, error) {
	line := shellQuote(append([]string{name}, args...))

	if r.DryRun {
		fmt.Printf("  + %s\n", line)
		return cmdResult{}, nil
	}

	fmt.Printf("  $ %s\n", line)

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	res := cmdResult{Stdout: stdout.String(), Stderr: stderr.String()}
	printIndented(res.Stdout)
	printIndented(res.Stderr)

	if err != nil {
		return res, fmt.Errorf("%s failed: %w", name, err)
	}
	return res, nil
}

func printIndented(s string) {
	s = strings.TrimRight(s, "\n")
	if s == "" {
		return
	}
	for _, ln := range strings.Split(s, "\n") {
		fmt.Printf("    | %s\n", ln)
	}
}

// shellQuote renders an argv as a single copy-pasteable shell line.
func shellQuote(argv []string) string {
	parts := make([]string, len(argv))
	for i, a := range argv {
		parts[i] = quoteArg(a)
	}
	return strings.Join(parts, " ")
}

func quoteArg(s string) string {
	if s == "" {
		return "''"
	}
	if !strings.ContainsAny(s, " \t\n'\"\\$`*?{}[]()&;|<>#~") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
