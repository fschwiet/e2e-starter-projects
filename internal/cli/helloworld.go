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
