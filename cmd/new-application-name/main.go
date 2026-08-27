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
