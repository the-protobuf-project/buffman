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

// flags holds the command-line flags for the `convert flatbuffers` command.
var (

	// flatbufferDir specifies the output directory for generated FlatBuffer files.
	flatbufferDir string
)

// flatbuffersCmd represents the cobra command to convert Protocol Buffer (.proto)
// files to FlatBuffer (.fbs) schema files.
var flatbuffersCmd = &cobra.Command{
	Use:   "flatbuffers",
	Short: "Converts Protocol Buffer (.proto) files to FlatBuffer (.fbs) schema files",
	Long: `The flatbuffers command parses a directory of Protocol Buffer (.proto) files
and converts them into FlatBuffer (.fbs) schema files. This facilitates
migration and interoperability between systems using these two serialization formats.

The command requires an input directory containing the source .proto files. An
output directory can be specified; if omitted, the generated .fbs files will be
placed in the current working directory.

Example:
  buffman convert flatbuffers --proto_dir=./path/to/protos --output_dir=./gen/fbs`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := parser.NewParser(parser.Flatbuffers)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := c.Parse(context.Background(), options.ParseOptions{
			InputDir:  protoDir,
			OutputDir: flatbufferDir,
		}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Successfully converted .proto files in %s to .fbs files in %s\n", protoDir, flatbufferDir)
	},
}

// init registers and configures the flags for the flatbuffersCmd.
func init() {
	// Bind the protoDir and flatbufferDir variables to command-line flags.
	flatbuffersCmd.Flags().StringVarP(&protoDir, "proto_dir", "I", "", "Directory containing the source .proto files")
	flatbuffersCmd.Flags().StringVarP(&flatbufferDir, "output_dir", "o", "./", "Output directory for the generated FlatBuffer (.fbs) files")

	// Mark the proto_dir flag as mandatory for command execution.
	flatbuffersCmd.MarkFlagRequired("proto_dir")
}
