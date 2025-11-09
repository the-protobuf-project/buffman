// package convert implements the "convert" command and its subcommands
// for the Buffman CLI. It provides functionality for converting between
// different buffer schema formats.
package convert

import "github.com/spf13/cobra"

// protoDir specifies the input directory containing .proto files.
var protoDir string

// ConvertCmd represents the base "convert" command.
// This command acts as a parent for all schema conversion subcommands,
// such as converting Protocol Buffers to FlatBuffers.
var ConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Groups subcommands for converting between schema formats",
	Long: `The 'convert' command provides a suite of tools for converting buffer schemas
from one format to another.

For example, you can use its subcommands to convert Protocol Buffers (.proto)
files into FlatBuffers (.fbs) schemas. Use a specific subcommand to define
the target output format.`,
	// A Run function is not needed here as this command only serves as a
	// gateway to its subcommands.
}

// The init function is used to attach subcommands to the ConvertCmd.
// This ensures that commands like 'convert flatbuffers' are available.
func init() {
	// flatbuffersCmd is defined in another file within this package (e.g., flatbuffers.go)
	ConvertCmd.AddCommand(flatbuffersCmd, nanobuffersCmd)
}
