// Package completion provides a custom implementation of the Cobra completion command.
// This package overrides Cobra's default completion command, preventing it from
// appearing in the list of available subcommands.
package completion

import "github.com/spf13/cobra"

// CompletionCmd defines the command for generating shell autocompletion scripts.
// The command is hidden by default in the init function to prevent it
// from being listed as a discoverable subcommand.
var CompletionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate the autocompletion script for the specified shell",
}

func init() {
	CompletionCmd.Hidden = true
}
