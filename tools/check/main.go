// Command check runs the verification pipeline.
//
// Usage:
//
//	go run ./tools/check          # every step, in order, stopping at the first failure
//	go run ./tools/check <step>   # a single step
//
// Step 1 (format) rewrites source files. See ADR-0003.
//
// This is a plain Go program rather than a Makefile or Taskfile so that the Go
// toolchain remains the only prerequisite. See ADR-0002.
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// binaryName is the output name of the built CLI. Rename it alongside the module
// path when starting a new app.
const binaryName = "new-application-name"

// step is one stage of the verification pipeline.
type step struct {
	name string
	args []string
}

// golangciLint builds an argument list that invokes the linter through its pinned
// tool modfile, so no global install is required. See ADR-0001.
func golangciLint(args ...string) []string {
	base := []string{"go", "tool", "-modfile=golangci-lint.mod", "golangci-lint"}

	return append(base, args...)
}

// outputPath is where the build step writes the binary.
func outputPath() string {
	name := binaryName
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	return "bin/" + name
}

// pipeline returns the steps in execution order: cheapest first, with the one
// hard constraint that build precedes e2e.
func pipeline() []step {
	return []step{
		{name: "format", args: golangciLint("fmt", "./...")},
		{name: "lint", args: golangciLint("run", "./...")},
		{name: "unit", args: []string{"go", "test", "./..."}},
		{name: "build", args: []string{"go", "build", "-o", outputPath(), "./cmd/" + binaryName}},
		{name: "e2e", args: []string{"go", "test", "-tags", "e2e", "./test/e2e/..."}},
	}
}

// stepByName finds a step by name.
func stepByName(steps []step, name string) (step, bool) {
	for _, s := range steps {
		if s.name == name {
			return s, true
		}
	}

	return step{}, false
}

// stepNames lists the step names, for error messages.
func stepNames(steps []step) []string {
	names := make([]string, 0, len(steps))
	for _, s := range steps {
		names = append(names, s.name)
	}

	return names
}

// selectSteps decides which steps to run for the given command-line arguments.
func selectSteps(args []string) ([]step, error) {
	all := pipeline()

	if len(args) == 0 {
		return all, nil
	}

	if len(args) > 1 {
		return nil, fmt.Errorf("expected at most one step name, got %d", len(args))
	}

	name := args[0]

	s, ok := stepByName(all, name)
	if !ok {
		return nil, fmt.Errorf("unknown step %q; valid steps are: %s", name, strings.Join(stepNames(all), ", "))
	}

	// e2e runs the built binary, so rebuild first. Otherwise a stale binary in
	// bin/ could produce a false pass for source that no longer compiles.
	if name == "e2e" {
		build, found := stepByName(all, "build")
		if !found {
			return nil, errors.New("internal error: build step missing from pipeline")
		}

		return []step{build, s}, nil
	}

	return []step{s}, nil
}

// runStep executes one step, streaming its output so long steps show progress.
func runStep(s step) error {
	fmt.Printf("==> %s\n", s.name)

	cmd := exec.Command(s.args[0], s.args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("step %q failed: %w", s.name, err)
	}

	return nil
}

func main() {
	steps, err := selectSteps(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "check:", err)
		os.Exit(2)
	}

	for _, s := range steps {
		if err := runStep(s); err != nil {
			fmt.Fprintln(os.Stderr, "check:", err)
			os.Exit(1)
		}
	}
}
