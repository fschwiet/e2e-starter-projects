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
