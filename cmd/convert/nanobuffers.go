// Package convert provides commands for converting between different data serialization schema formats.
package convert

import (
	"context"
	"fmt"
	"os"

	"github.com/machanirobotics/buffman/internal/options"
	"github.com/machanirobotics/buffman/internal/parser"
	"github.com/spf13/cobra"
)

// flags holds the command-line flags for the `convert nanobuffers` command.
var (
	// nanobufferDir specifies the output directory for generated Nanobuffer files.
	nanobufferDir string
)

// nanobuffersCmd represents the cobra command to convert Protocol Buffer (.proto)
// files to Nanobuffer (.nanobuf) schema files.
var nanobuffersCmd = &cobra.Command{
	Use:   "nanobuffers",
	Short: "Converts Protocol Buffer (.proto) files to Nanobuffer (.nanobuf) schema files",
	Long: `The nanobuffers command parses a directory of Protocol Buffer (.proto) files
and converts them into Nanobuffer (.nanobuf) schema files. This facilitates
migration and interoperability between systems using these two serialization formats.

The command requires an input directory containing the source .proto files. An
output directory can be specified; if omitted, the generated .nanobuf files will be
placed in the current working directory.

Example:
  buffman convert nanobuffers --proto_dir=./path/to/protos --output_dir=./gen/nanobuf`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := parser.NewParser(parser.Nanobuffers)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := c.Parse(context.Background(), options.ParseOptions{
			InputDir:  protoDir,
			OutputDir: nanobufferDir,
		}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Successfully converted .proto files in %s to .nanobuf files in %s\n", protoDir, nanobufferDir)
	},
}

// init registers and configures the flags for the nanobuffersCmd.
func init() {
	// Bind the protoDir and nanobufferDir variables to command-line flags.
	nanobuffersCmd.Flags().StringVarP(&protoDir, "proto_dir", "I", "", "Directory containing the source .proto files")
	nanobuffersCmd.Flags().StringVarP(&nanobufferDir, "output_dir", "o", "./", "Output directory for the generated Nanobuffer (.nanobuf) files")

	// Mark the proto_dir flag as mandatory for command execution.
	nanobuffersCmd.MarkFlagRequired("proto_dir")
}
