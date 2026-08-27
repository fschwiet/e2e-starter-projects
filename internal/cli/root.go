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
