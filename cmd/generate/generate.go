// Package generate provides the primary "generate" command, which serves as an
// entry point for all code generation subcommands within the Buffman CLI.
package generate

import (
	"context"
	"fmt"
	"os"

	"github.com/machanirobotics/buffman/internal/runner"
	"github.com/spf13/cobra"
)

// buffmanConfigPath holds the path to the buffman.yaml configuration file,
// provided via the -f or --file flag.
var (
	buffmanConfigPath string
)

// GenerateCmd represents the base "generate" command.
// It acts as a dispatcher, either executing a generation run based on a
// configuration file or delegating to a specific code generation subcommand
// like 'flatbuffers'.
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates code from schema files using a config or subcommands",
	Long: `The 'generate' command is the entry point for all code generators.

It can be used in two main ways:

1. With a configuration file:
   Provide a path to a buffman.yaml file using the --file flag. This is the
   recommended approach for managing all your code generation settings in one place.

   Example:
     buffman generate --file buffman.yaml

2. With a subcommand:
   Invoke a specific generator directly, providing its options via flags.

   Example:
     buffman generate flatbuffers --language=go --target_dir=./gen`,
	// Run executes the logic for the generate command. If a configuration file
	// path is provided, it uses the runner. Otherwise, it prints the help text.
	Run: func(cmd *cobra.Command, args []string) {
		if buffmanConfigPath != "" {
			handleWithConfig(buffmanConfigPath)
		} else {
			cmd.Help()
		}
	},
}

// init registers flags and adds subcommands to GenerateCmd.
func init() {
	// The --file flag allows users to specify a configuration file for generation.
	GenerateCmd.Flags().StringVarP(&buffmanConfigPath, "file", "f", "", "Path to the buffman.yaml configuration file")
	// Add subcommands for specific generators.
	GenerateCmd.AddCommand(flatbuffersCmd)
}

// handleWithConfig initializes and executes a runner to perform code generation
// based on the settings in the provided configuration file. It exits the program
// if the runner encounters an error.
func handleWithConfig(configPath string) {
	run := runner.NewRunner()
	if err := run.Run(context.Background(), configPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
