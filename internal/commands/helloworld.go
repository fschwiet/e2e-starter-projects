// Package commands holds the logic behind each CLI command. Functions here are
// pure: no cobra, no I/O, no global state. That is what makes them unit-testable
// without constructing a command.
package commands

// Greet returns the greeting printed by the hello-world command.
func Greet() string {
	return "Hello, world!"
}
