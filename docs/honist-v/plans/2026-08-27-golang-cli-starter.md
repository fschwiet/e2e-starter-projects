# Go CLI Starter Kit Implementation Plan

**Goal:** Build a minimal "hello, world" Go CLI starter kit on the `golang-cli` branch with a local verification pipeline covering formatting, linting, unit tests and end-to-end tests.

**Architecture:** A thin `main` package delegates to cobra wiring in `internal/cli`, which delegates to pure logic in `internal/commands`. That split is the pattern the starter kit exists to demonstrate — cobra wiring never contains logic, so logic is unit-testable without constructing a cobra command. A Go program at `tools/check` runs the five pipeline steps by shelling out, so the pipeline needs no task runner installed.

**Tech Stack:** Go 1.26.7, `spf13/cobra` v1.10.2, golangci-lint v2.13.1 (pinned as a Go tool dependency in a dedicated modfile), gofumpt + goimports (as golangci-lint formatters), stdlib `testing`.

## Global Constraints

Every task's requirements implicitly include this section.

- **`go` is the only prerequisite.** Never add a step that requires installing another tool. golangci-lint is compiled on demand from `golangci-lint.mod`.
- Module path is exactly `github.com/your-org/new-application-name`.
- The `go` directive in `go.mod` is `go 1.26.7`.
- cobra is pinned at exactly `v1.10.2`.
- golangci-lint is pinned at exactly `v2.13.1`.
- The placeholder application name is `new-application-name` everywhere. The example subcommand is `hello-world` and is independent of it.
- `tools/check` uses **only the standard library** (`fmt`, `os`, `os/exec`, `runtime`, `strings`, `errors`). It must add nothing to `go.mod`.
- No `-race` anywhere in the pipeline.
- No CI workflow, no release automation, no `CONTEXT.md`.
- Unit tests live beside the code they test. Only e2e tests get a separate directory.
- Commit at the end of every task.

### Environment notes

These are verified facts about the toolchain. They will otherwise cost the implementer an hour each.

1. **PowerShell splits an unquoted `-modfile=golangci-lint.mod` argument at the extension.** It surfaces as `file does not have .mod extension`, or as `'go mod init' accepts at most one argument`. In PowerShell, quote the flag: `go tool "-modfile=golangci-lint.mod" golangci-lint run ./...`. In bash/git-bash, no quoting is needed. `tools/check` shells out via `os/exec` and is immune either way.

2. **`go get -tool` must be given the `/cmd/golangci-lint` package path, not the module root.** The module-root form prints `cannot find module providing package …` but **still exits 0**, writing a `tool` directive pointing at a non-main package. It only fails later, at `go tool` time, with the misleading `no required module provides package github.com/golangci/golangci-lint/v2`.

3. **The first golangci-lint run compiles it from source and takes minutes.** It is not hung. Later runs hit Go's build cache.

4. **A package whose files are all excluded by a build tag is not an error** under a `./...` pattern. `go test ./...` and `go build ./...` skip `test/e2e/` silently and exit 0.

---

## File Structure

| File | Responsibility |
| --- | --- |
| `go.mod` / `go.sum` | Application module. Depends only on cobra. |
| `golangci-lint.mod` / `.sum` | Pinned linter tool dependency, isolated from the app's dependency graph. |
| `.golangci.yml` | Formatter and linter configuration; exposes the `e2e` build tag to the linter. |
| `.gitignore` | Ignores build output. |
| `cmd/new-application-name/main.go` | Process entry point. Calls `cli.Execute()`, maps error to exit code 1. Nothing else. |
| `internal/cli/root.go` | Root cobra command, version var, `Execute()`. |
| `internal/cli/helloworld.go` | `hello-world` subcommand wiring only. |
| `internal/commands/helloworld.go` | `Greet() string` — pure logic, no cobra, no I/O. |
| `internal/commands/helloworld_test.go` | Unit test for `Greet()`. |
| `tools/check/main.go` | Verification pipeline runner. |
| `tools/check/main_test.go` | Unit tests for the runner's step-selection logic. |
| `test/e2e/cli_test.go` | E2E tests against the built binary. Tagged `//go:build e2e`. |
| `README.md` | Human-facing docs. |
| `AGENTS.md` | Agent-facing pipeline docs. |
| `CLAUDE.md` | One line: `@AGENTS.md`. |

### Deliberate addition beyond the spec

The spec's Testing section names only `internal/commands/helloworld_test.go` as a unit test. This plan adds **`tools/check/main_test.go`**, covering `selectSteps`. The runner contains the only non-trivial branching logic in the kit (the `e2e`-implies-`build` rule and the unknown-step error), and that logic is otherwise untested by anything. It is a pure function, so the test is cheap. Flagging it as a conscious addition rather than scope creep.

---

## Task 1: Module bootstrap and the greeting function

Establishes the module and the pure-logic package, test-first.

**Files:**

- Create: `go.mod`, `.gitignore`
- Create: `internal/commands/helloworld.go`
- Test: `internal/commands/helloworld_test.go`

**Interfaces:**

- Consumes: nothing.
- Produces: package `github.com/your-org/new-application-name/internal/commands`, exporting `func Greet() string`, which returns `"Hello, world!"`.

- [ ] **Step 1: Initialize the module**

Run from the repo root:

```bash
go mod init github.com/your-org/new-application-name
```

Expected: `go: creating new go.mod: module github.com/your-org/new-application-name`

Verify the `go` directive says `go 1.26.7`:

```bash
cat go.mod
```

- [ ] **Step 2: Create `.gitignore`**

Create `.gitignore`:

```gitignore
# Build output from `go run ./tools/check build`
/bin/

# Go build and test artifacts
*.exe
*.test
*.out
```

- [ ] **Step 3: Write the failing test**

Create `internal/commands/helloworld_test.go`:

```go
package commands

import "testing"

func TestGreet(t *testing.T) {
	got := Greet()
	want := "Hello, world!"

	if got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}
```

- [ ] **Step 4: Run the test to verify it fails**

```bash
go test ./internal/commands/
```

Expected: FAIL — `undefined: Greet`

- [ ] **Step 5: Write the minimal implementation**

Create `internal/commands/helloworld.go`:

```go
// Package commands holds the logic behind each CLI command. Functions here are
// pure: no cobra, no I/O, no global state. That is what makes them unit-testable
// without constructing a command, and it is the pattern every new command in this
// starter kit should follow.
package commands

// Greet returns the greeting printed by the hello-world command.
func Greet() string {
	return "Hello, world!"
}
```

- [ ] **Step 6: Run the test to verify it passes**

```bash
go test ./internal/commands/
```

Expected: `ok  	github.com/your-org/new-application-name/internal/commands`

- [ ] **Step 7: Commit**

```bash
git add go.mod .gitignore internal/commands/
git commit -m "feat: add module and Greet() with unit test"
```

---

## Task 2: Cobra CLI wiring and entry point

Wires `Greet()` up to a real command. Verified by running the CLI; locked in by e2e tests in Task 3.

**Files:**

- Create: `internal/cli/root.go`
- Create: `internal/cli/helloworld.go`
- Create: `cmd/new-application-name/main.go`
- Modify: `go.mod`, `go.sum` (adds cobra)

**Interfaces:**

- Consumes: `commands.Greet() string` from Task 1.
- Produces:
  - `func cli.Execute() error` — runs the CLI, returns an error if the command failed. cobra has already written the message to stderr, so callers must not print it again.
  - unexported `func newRootCommand() *cobra.Command`, `func newHelloWorldCommand() *cobra.Command`.
  - binary behaviour later relied on by Task 3's e2e tests: `hello-world` prints `Hello, world!\n` to stdout with exit 0; an unknown command exits 1 with `unknown command` on stderr and empty stdout.

- [ ] **Step 1: Add the cobra dependency**

```bash
go get github.com/spf13/cobra@v1.10.2
```

Expected: `go: added github.com/spf13/cobra v1.10.2`

- [ ] **Step 2: Write the root command**

Create `internal/cli/root.go`:

```go
// Package cli wires up the command-line interface. It contains no business
// logic: every command delegates to a function in internal/commands. Keeping
// this boundary is what lets the logic be tested without cobra.
package cli

import "github.com/spf13/cobra"

// version is reported by --version. It is a var rather than a const so a release
// build can override it without code changes:
//
//	go build -ldflags "-X github.com/your-org/new-application-name/internal/cli.version=1.2.3" ./cmd/new-application-name
var version = "0.0.1"

// newRootCommand builds the root command and registers every subcommand.
func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:     "new-application-name",
		Short:   "A starter kit for Go CLI applications",
		Version: version,
	}

	root.AddCommand(newHelloWorldCommand())

	return root
}

// Execute runs the CLI. On failure it returns an error; cobra has already
// written the message to stderr, so the caller must not print it again.
func Execute() error {
	return newRootCommand().Execute()
}
```

- [ ] **Step 3: Write the hello-world subcommand**

Create `internal/cli/helloworld.go`:

```go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/your-org/new-application-name/internal/commands"
)

// newHelloWorldCommand builds the example subcommand. Note how little it does:
// it calls one function and writes the result. Copy this shape for new commands.
func newHelloWorldCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "hello-world",
		Short: "Print a greeting",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Writing to cmd.OutOrStdout() rather than os.Stdout keeps the
			// command testable. The returned error is deliberately checked
			// rather than discarded: errcheck flags the idiomatic
			// fmt.Fprintln(...) call, and returning it is the honest fix.
			_, err := fmt.Fprintln(cmd.OutOrStdout(), commands.Greet())

			return err
		},
	}
}
```

- [ ] **Step 4: Write the entry point**

Create `cmd/new-application-name/main.go`:

```go
// Command new-application-name is the CLI entry point. It stays deliberately
// thin: all wiring lives in internal/cli, all logic in internal/commands.
package main

import (
	"os"

	"github.com/your-org/new-application-name/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		// cobra already printed the message to stderr.
		os.Exit(1)
	}
}
```

- [ ] **Step 5: Verify the CLI works**

```bash
go run ./cmd/new-application-name hello-world
```

Expected: `Hello, world!`

```bash
go run ./cmd/new-application-name --version
```

Expected: `new-application-name version 0.0.1`

```bash
go run ./cmd/new-application-name bogus; echo "exit: $?"
```

Expected: `Error: unknown command "bogus" for "new-application-name"`, then `Run 'new-application-name --help' for usage.`, then `exit: 1`

- [ ] **Step 6: Verify the existing unit test still passes**

```bash
go test ./...
```

Expected: `ok` for `internal/commands`, `[no test files]` for the others.

- [ ] **Step 7: Commit**

```bash
git add go.mod go.sum internal/cli/ cmd/
git commit -m "feat: add cobra CLI wiring and entry point"
```

---

## Task 3: End-to-end tests against the built binary

Adds the e2e layer and the build-tag isolation that keeps `go test ./...` fast.

**Files:**

- Create: `test/e2e/cli_test.go`

**Interfaces:**

- Consumes: the built binary at `bin/new-application-name` (`.exe` on Windows), and its behaviour as specified in Task 2's Produces block.
- Produces: nothing importable — this package is tagged `e2e` and is a leaf.

- [ ] **Step 1: Write the failing e2e tests**

Create `test/e2e/cli_test.go`:

```go
//go:build e2e

// Package e2e exercises the built binary as a subprocess, so it tests the real
// shipped artifact rather than a package.
//
// These two tests are exemplars, not coverage. Between them they show the two
// shapes an assertion takes — the success path via stdout, and the failure path
// via stderr and exit code. Copy whichever fits when adding a command.
//
// The e2e build tag keeps these out of `go test ./...`, which would otherwise
// fail confusingly on a clean checkout where the binary has not been built.
package e2e

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// binaryPath returns the binary produced by the pipeline's build step. Go runs
// tests with the working directory set to the package directory, so this
// relative path resolves deterministically.
func binaryPath(t *testing.T) string {
	t.Helper()

	name := "new-application-name"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	path := filepath.Join("..", "..", "bin", name)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("binary not found at %s: run \"go run ./tools/check build\" first", path)
	}

	return path
}

// runCLI executes the binary and returns its stdout, stderr and exit code.
func runCLI(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	cmd := exec.Command(binaryPath(t), args...)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err := cmd.Run()

	var exitErr *exec.ExitError
	switch {
	case err == nil:
		exitCode = 0
	case errors.As(err, &exitErr):
		exitCode = exitErr.ExitCode()
	default:
		t.Fatalf("running %s: %v", cmd, err)
	}

	return outBuf.String(), errBuf.String(), exitCode
}

func TestHelloWorldPrintsGreeting(t *testing.T) {
	stdout, stderr, exitCode := runCLI(t, "hello-world")

	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0 (stderr: %s)", exitCode, stderr)
	}

	if want := "Hello, world!\n"; stdout != want {
		t.Errorf("stdout = %q, want %q", stdout, want)
	}
}

func TestUnknownCommandFails(t *testing.T) {
	// Asserts a substring of the unknown-*command* message rather than cobra's
	// exact text. The unknown-command case prints a short help hint with no
	// usage block, so this assertion survives the addition of new commands and
	// flags; asserting the full usage block would not.
	stdout, stderr, exitCode := runCLI(t, "bogus")

	if exitCode != 1 {
		t.Errorf("exit code = %d, want 1", exitCode)
	}

	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("stderr = %q, want it to contain %q", stderr, "unknown command")
	}

	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
}
```

- [ ] **Step 2: Run the e2e tests to verify they fail**

The binary has not been built yet.

```bash
go test -tags e2e ./test/e2e/...
```

Expected: FAIL — `binary not found at ..\..\bin\new-application-name.exe: run "go run ./tools/check build" first`

This confirms the actionable error message works.

- [ ] **Step 3: Build the binary**

On bash / git-bash:

```bash
go build -o bin/new-application-name.exe ./cmd/new-application-name
```

On non-Windows, drop the `.exe`. Task 4 automates this; here it is manual so the tests can be proven.

- [ ] **Step 4: Run the e2e tests to verify they pass**

```bash
go test -tags e2e ./test/e2e/...
```

Expected: `ok  	github.com/your-org/new-application-name/test/e2e`

- [ ] **Step 5: Verify the build tag isolates the e2e package**

```bash
go test ./...
```

Expected: `ok` for `internal/commands`, `[no test files]` elsewhere, and **no mention of `test/e2e`** — and exit 0. A package with all files excluded by a build tag is skipped silently, not treated as an error.

- [ ] **Step 6: Commit**

```bash
git add test/e2e/
git commit -m "test: add e2e tests against the built binary"
```

---

## Task 4: Pin golangci-lint and configure lint + format

Sets up the linter with no global install, and configures it to see the tagged e2e files.

**Files:**

- Create: `golangci-lint.mod`, `golangci-lint.sum`
- Create: `.golangci.yml`

**Interfaces:**

- Consumes: nothing from earlier tasks.
- Produces: two invocations later used verbatim by Task 5's runner:
  - `go tool -modfile=golangci-lint.mod golangci-lint fmt ./...` — applies formatting, rewrites files
  - `go tool -modfile=golangci-lint.mod golangci-lint run ./...` — reports lint **and** formatter violations, exit 1 on findings

- [ ] **Step 1: Create the dedicated tool module**

```bash
go mod init -modfile=golangci-lint.mod github.com/your-org/new-application-name/golangci-lint
```

Expected: `go: creating new go.mod: module github.com/your-org/new-application-name/golangci-lint`

In PowerShell, quote the flag: `go mod init "-modfile=golangci-lint.mod" github.com/your-org/new-application-name/golangci-lint`

- [ ] **Step 2: Pin golangci-lint as a tool dependency**

The package path must end in `/cmd/golangci-lint`. See Environment note 2 — the module-root form fails silently.

```bash
go get -tool -modfile=golangci-lint.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.1
```

Expected: a long list of `go: added …` lines, exit 0, and **no** `cannot find module providing package` line.

- [ ] **Step 3: Verify the tool directive is correct**

```bash
grep tool golangci-lint.mod
```

Expected: `tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint`

If it reads `tool github.com/golangci/golangci-lint/v2` (no `/cmd/...`), Step 2 used the wrong path. Delete `golangci-lint.mod` and `golangci-lint.sum` and redo Steps 1–2.

- [ ] **Step 4: Confirm `golangci-lint.sum` exists**

```bash
ls golangci-lint.mod golangci-lint.sum
```

Both must exist and both get committed.

- [ ] **Step 5: Write the configuration**

Create `.golangci.yml`:

```yaml
version: "2"

# The e2e tests are behind a build tag. Without this, golangci-lint cannot see
# test/e2e/ at all and that code silently rots unlinted.
run:
  build-tags:
    - e2e

# gofumpt is a stricter superset of gofmt; goimports handles import grouping and
# pruning, which gofumpt does not do.
formatters:
  enable:
    - gofumpt
    - goimports

# Linters are golangci-lint's defaults: errcheck, govet, ineffassign, staticcheck,
# unused. Deliberately not expanded — add linters as the project matures rather
# than starting by deleting ones you did not choose.
```

- [ ] **Step 6: Run the formatter**

This is the first run and compiles golangci-lint from source. Expect minutes. It is not hung.

```bash
go tool -modfile=golangci-lint.mod golangci-lint fmt ./...
```

Expected: exit 0. It may rewrite files; that is normal.

- [ ] **Step 7: Run the linter and fix anything it reports**

```bash
go tool -modfile=golangci-lint.mod golangci-lint run ./...
```

Expected: `0 issues.`, exit 0.

If `errcheck` flags an unchecked error, fix the code by handling the error — do **not** add an exclusion to `.golangci.yml`.

- [ ] **Step 8: Verify the linter sees the tagged e2e package**

Temporarily append an unused function to `test/e2e/cli_test.go`:

```go
func unusedProbe() string { return "x" }
```

Then:

```bash
go tool -modfile=golangci-lint.mod golangci-lint run ./...
```

Expected: `test/e2e/cli_test.go:…: func unusedProbe is unused (unused)`

This proves `run.build-tags` is working. **Delete `unusedProbe` afterwards** and re-run to confirm `0 issues.`

- [ ] **Step 9: Commit**

```bash
git add golangci-lint.mod golangci-lint.sum .golangci.yml
git commit -m "build: pin golangci-lint v2.13.1 and configure lint and format"
```

---

## Task 5: The verification pipeline runner

Replaces the manual commands from Tasks 3 and 4 with one entry point.

**Files:**

- Create: `tools/check/main.go`
- Test: `tools/check/main_test.go`

**Interfaces:**

- Consumes: the golangci-lint invocations from Task 4; the build and test commands used manually in Task 3.
- Produces:
  - `go run ./tools/check` — runs all five steps in order, stopping at the first failure.
  - `go run ./tools/check <step>` — runs one step; `e2e` runs `build` first.
  - internal, exercised by tests: `type step struct { name string; args []string }`, `func pipeline() []step`, `func stepByName(steps []step, name string) (step, bool)`, `func selectSteps(args []string) ([]step, error)`.

- [ ] **Step 1: Write the failing tests**

Create `tools/check/main_test.go`:

```go
package main

import (
	"strings"
	"testing"
)

func stepNamesOf(steps []step) []string {
	names := make([]string, 0, len(steps))
	for _, s := range steps {
		names = append(names, s.name)
	}

	return names
}

func TestSelectStepsWithNoArgsRunsWholePipeline(t *testing.T) {
	steps, err := selectSteps(nil)
	if err != nil {
		t.Fatalf("selectSteps(nil) returned error: %v", err)
	}

	got := strings.Join(stepNamesOf(steps), ",")
	want := "format,lint,unit,build,e2e"

	if got != want {
		t.Errorf("step names = %q, want %q", got, want)
	}
}

func TestSelectStepsWithOneNameRunsOnlyThatStep(t *testing.T) {
	steps, err := selectSteps([]string{"lint"})
	if err != nil {
		t.Fatalf("selectSteps([lint]) returned error: %v", err)
	}

	got := strings.Join(stepNamesOf(steps), ",")
	if want := "lint"; got != want {
		t.Errorf("step names = %q, want %q", got, want)
	}
}

func TestSelectStepsBuildsBeforeE2E(t *testing.T) {
	// Running e2e alone against a stale binary could report a false pass, so
	// the runner always rebuilds first.
	steps, err := selectSteps([]string{"e2e"})
	if err != nil {
		t.Fatalf("selectSteps([e2e]) returned error: %v", err)
	}

	got := strings.Join(stepNamesOf(steps), ",")
	if want := "build,e2e"; got != want {
		t.Errorf("step names = %q, want %q", got, want)
	}
}

func TestSelectStepsRejectsUnknownStep(t *testing.T) {
	_, err := selectSteps([]string{"nope"})
	if err == nil {
		t.Fatal("selectSteps([nope]) returned nil error, want an error")
	}

	// The message must list the valid names so the user can recover.
	for _, name := range []string{"format", "lint", "unit", "build", "e2e"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not mention valid step %q", err, name)
		}
	}
}

func TestSelectStepsRejectsMultipleArgs(t *testing.T) {
	if _, err := selectSteps([]string{"lint", "unit"}); err == nil {
		t.Fatal("selectSteps([lint unit]) returned nil error, want an error")
	}
}
```

- [ ] **Step 2: Run the tests to verify they fail**

```bash
go test ./tools/check/
```

Expected: FAIL — `undefined: step`, `undefined: selectSteps`

- [ ] **Step 3: Write the runner**

Create `tools/check/main.go`:

```go
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
```

- [ ] **Step 4: Run the tests to verify they pass**

```bash
go test ./tools/check/
```

Expected: `ok  	github.com/your-org/new-application-name/tools/check`

- [ ] **Step 5: Run the whole pipeline**

```bash
go run ./tools/check
```

Expected, in order:

```
==> format
==> lint
0 issues.
==> unit
ok  	github.com/your-org/new-application-name/internal/commands
ok  	github.com/your-org/new-application-name/tools/check
==> build
==> e2e
ok  	github.com/your-org/new-application-name/test/e2e
```

The `unit` step also prints `[no test files]` lines for `cmd/new-application-name`,
`internal/cli` and `test/e2e`. That is expected, not a problem.

Exit code 0.

- [ ] **Step 6: Verify single-step dispatch and the e2e-builds-first rule**

```bash
go run ./tools/check lint
```

Expected: only `==> lint`.

```bash
go run ./tools/check e2e
```

Expected: `==> build` then `==> e2e`.

```bash
go run ./tools/check nope; echo "exit: $?"
```

Expected: `check: unknown step "nope"; valid steps are: format, lint, unit, build, e2e` and `exit: 2`

- [ ] **Step 7: Confirm the runner added no dependencies**

```bash
cat go.mod
```

Expected: `github.com/spf13/cobra v1.10.2` is the only direct requirement. `tools/check`
imports nothing outside the standard library, so it must not have added anything.

- [ ] **Step 8: Commit**

```bash
git add tools/check/
git commit -m "build: add verification pipeline runner"
```

---

## Task 6: Documentation

The starter kit is only usable if its conventions are discoverable — by humans and by agents.

**Files:**

- Create: `README.md` (this replaces the repo-root `README.md` inherited from `main`)
- Create: `AGENTS.md`
- Create: `CLAUDE.md`

**Interfaces:**

- Consumes: every command established in Tasks 1–5.
- Produces: nothing importable.

- [ ] **Step 1: Write the README**

Overwrite `README.md`:

````markdown
# new-application-name

A starter kit for Go CLI applications: cobra, golangci-lint, gofumpt, unit tests and
end-to-end tests, with a local verification pipeline that needs nothing installed but Go.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26 or newer.

That is the entire list. golangci-lint is pinned in `golangci-lint.mod` and compiled on
demand by the Go toolchain — there is nothing else to install.

## Verification pipeline

Run everything:

```
go run ./tools/check
```

Run a single step:

```
go run ./tools/check lint
```

| # | Step | What it does | Rewrites files |
| - | ---- | ------------ | -------------- |
| 1 | `format` | applies gofumpt + goimports | **yes** |
| 2 | `lint` | golangci-lint: errcheck, govet, ineffassign, staticcheck, unused | no |
| 3 | `unit` | `go test ./...` | no |
| 4 | `build` | compiles to `bin/new-application-name` | no |
| 5 | `e2e` | runs the built binary as a subprocess | no |

Steps run in order and stop at the first failure.

**Step 1 rewrites your source files.** Formatting is applied rather than reported, so you
never have to run a formatter and then re-run the pipeline. The trade-off is that a
passing run can leave you with unstaged changes — "did it pass" and "is my tree clean"
are separate questions here.

**The first run is slow** — a minute or more — because the Go toolchain compiles
golangci-lint from source. Later runs hit the build cache and are quick. It is not hung.

`go run ./tools/check e2e` builds before testing, so it can never pass against a stale
binary.

`go test ./...` never runs the e2e tests; they are behind the `e2e` build tag.

## Running the CLI

```
go run ./cmd/new-application-name hello-world
# Hello, world!
```

Or build first:

```
go run ./tools/check build
./bin/new-application-name hello-world
```

## Project structure

```
cmd/new-application-name/   the main package — thin, just an entry point
internal/cli/               cobra wiring only, no logic
internal/commands/          the logic, as pure functions — unit-tested here
test/e2e/                   tests the built binary (build tag: e2e)
tools/check/                the verification pipeline runner
```

The important rule: **keep logic out of `internal/cli`.** Commands delegate to
`internal/commands`, which is what makes the logic testable without constructing a cobra
command. Follow that split when adding commands.

## Starting a new app

Find-and-replace `new-application-name` across the whole repo, then rename the
`cmd/new-application-name/` directory to match. The name appears in:

- `go.mod` — the module path
- `golangci-lint.mod` — the module path
- `cmd/new-application-name/` — the directory name
- `cmd/new-application-name/main.go` — the import of `internal/cli`
- `internal/cli/root.go` — the root command's `Use` field
- `internal/cli/helloworld.go` — the import of `internal/commands`
- `tools/check/main.go` — the `binaryName` constant
- `test/e2e/cli_test.go` — the binary name it executes
- `README.md` — the title and examples

The import paths matter: renaming the module without updating them leaves the project
uncompilable. A repo-wide replace covers all of the above at once.

`hello-world` is the example command name and is independent of the application name.
Replace it when you add real commands.

## Versioning

`--version` reads `version` in `internal/cli/root.go`, which defaults to `0.0.1`. It is a
`var`, not a `const`, so a release build can override it without touching the source:

```
go build -ldflags "-X github.com/your-org/new-application-name/internal/cli.version=1.2.3" ./cmd/new-application-name
```

The pipeline does not set it. There is deliberately no e2e test asserting the version —
that test would fail on every version bump.

## Adding the race detector

Not enabled by default: on `windows/amd64` it needs `CGO_ENABLED=1` and a C toolchain,
which would break the "only Go required" promise, and it is unsupported on
`windows/arm64`. To enable it, change the `unit` step in `tools/check/main.go` to:

```go
{name: "unit", args: []string{"go", "test", "-race", "./..."}},
```

## Running golangci-lint directly

```
go tool -modfile=golangci-lint.mod golangci-lint run ./...
```

In PowerShell the flag must be quoted, or the `.mod` extension gets split off the
argument:

```powershell
go tool "-modfile=golangci-lint.mod" golangci-lint run ./...
```
````

- [ ] **Step 2: Write AGENTS.md**

Create `AGENTS.md`:

```markdown
# Working in this repo

A Go CLI starter kit. `go` is the only prerequisite — golangci-lint is pinned in
`golangci-lint.mod` and compiled on demand. Do not install tooling.

## Verify a change

Run this before considering any change done:

```
go run ./tools/check
```

Steps run in order and stop at the first failure:

| # | Step | What it checks |
| - | ---- | -------------- |
| 1 | `format` | applies gofumpt + goimports — **rewrites source files** |
| 2 | `lint` | golangci-lint: errcheck, govet, ineffassign, staticcheck, unused |
| 3 | `unit` | `go test ./...` |
| 4 | `build` | compiles the binary into `bin/` |
| 5 | `e2e` | runs the built binary as a subprocess |

Run one step with `go run ./tools/check <step>`. `e2e` rebuilds first, so it never tests
a stale binary.

**`go run ./tools/check` mutates the working tree** — step 1 applies formatting. Unstaged
changes after a run are expected, not a failure.

`tools/check` is the entry point. Do not invent an ad-hoc pipeline.

The first run compiles golangci-lint from source and takes minutes. It is not hung.

## Where code goes

- `internal/commands/` — logic, as pure functions. No cobra, no I/O. Unit-tested here,
  in `*_test.go` files beside the code.
- `internal/cli/` — cobra wiring only. Each command calls one function and writes the
  result.
- `cmd/new-application-name/` — the `main` package. Thin.
- `test/e2e/` — tests the built binary. Tagged `//go:build e2e`.
- `tools/check/` — the pipeline runner. Standard library only; it must add nothing to
  `go.mod`.

Keep logic out of `internal/cli`. That split is what makes commands testable without
constructing a cobra command, and it is the pattern this starter kit exists to show.

## Gotchas

- `go test ./...` does not run the e2e tests — they are behind the `e2e` build tag. Use
  `go run ./tools/check e2e`.
- Invoking golangci-lint directly from PowerShell needs the flag quoted, or the `.mod`
  extension is split off the argument:
  `go tool "-modfile=golangci-lint.mod" golangci-lint run ./...`
- If errcheck flags an unchecked error, handle the error. Do not add an exclusion to
  `.golangci.yml`.

## Design decisions

`docs/adr/` records why the non-obvious choices were made — the pinned tool modfile, the
Go pipeline runner, and the mutating format step. Read before changing any of them.
```

- [ ] **Step 3: Write CLAUDE.md**

Create `CLAUDE.md` with exactly one line:

```markdown
@AGENTS.md
```

- [ ] **Step 4: Verify the pipeline still passes**

```bash
go run ./tools/check
```

Expected: all five steps, exit 0.

- [ ] **Step 5: Verify the documented commands actually work**

Run each command quoted in the README and confirm it behaves as documented:

```bash
go run ./cmd/new-application-name hello-world
go run ./tools/check build && ./bin/new-application-name hello-world
go tool -modfile=golangci-lint.mod golangci-lint run ./...
```

Expected: `Hello, world!`, `Hello, world!`, `0 issues.`

- [ ] **Step 6: Commit**

```bash
git add README.md AGENTS.md CLAUDE.md
git commit -m "docs: add README, AGENTS.md and CLAUDE.md"
```

---

## Task 7: Final verification from a clean state

Proves the starter kit's central promise — clone it, run one command, everything works — which no earlier task tests, because each ran with warm caches and a dirty tree.

**Files:** none created or modified.

**Interfaces:** none.

- [ ] **Step 1: Confirm the working tree is clean**

```bash
git status --short
```

Expected: no output. If `bin/` appears, `.gitignore` from Task 1 is wrong.

- [ ] **Step 2: Verify the committed file set**

```bash
git ls-files
```

Expected exactly this list:

```
.gitignore
.golangci.yml
AGENTS.md
CLAUDE.md
README.md
cmd/new-application-name/main.go
docs/adr/0001-pin-golangci-lint-via-tool-modfile.md
docs/adr/0002-go-program-as-pipeline-runner.md
docs/adr/0003-format-mutates-inside-the-pipeline.md
docs/honist-v/plans/2026-08-27-golang-cli-starter.md
docs/honist-v/specs/2026-08-27-golang-cli-starter-design.md
go.mod
go.sum
golangci-lint.mod
golangci-lint.sum
internal/cli/helloworld.go
internal/cli/root.go
internal/commands/helloworld.go
internal/commands/helloworld_test.go
test/e2e/cli_test.go
tools/check/main.go
tools/check/main_test.go
```

Both `golangci-lint.mod` and `golangci-lint.sum` must be present.

- [ ] **Step 3: Verify a clean clone passes**

Clone into a directory outside the repo. Adjust the paths if your checkout differs:

```bash
VERIFY_DIR="$HOME/golang-cli-verify"
rm -rf "$VERIFY_DIR"
git clone -b golang-cli "C:/code/e2e-starter-kits" "$VERIFY_DIR"
cd "$VERIFY_DIR" && go run ./tools/check
```

Expected: all five steps, exit 0. This is the promise the starter kit makes — a fresh
checkout with only Go installed runs the whole pipeline.

- [ ] **Step 4: Confirm the clean clone left no unexpected changes**

```bash
git status --short
```

Expected: no output. If `format` rewrote files in a clean clone, formatting was not
applied before committing — go back, run `go run ./tools/check`, and commit the result.

- [ ] **Step 5: Clean up**

```bash
cd "C:/code/e2e-starter-kits" && rm -rf "$HOME/golang-cli-verify"
```

- [ ] **Step 6: Commit any fixes**

Only if Steps 1–4 surfaced problems:

```bash
git add -A
git commit -m "fix: address issues found in clean-clone verification"
```
