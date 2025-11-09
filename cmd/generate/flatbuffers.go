// Package generate provides commands for generating language-specific source code
// from various schema definition files.
package generate

import (
	"context"
	"fmt"
	"os"

	"github.com/machanirobotics/buffman/internal/configuration"
	"github.com/machanirobotics/buffman/internal/generate"
	"github.com/machanirobotics/buffman/internal/generate/language"
	"github.com/machanirobotics/buffman/internal/options"
	"github.com/spf13/cobra"
)

// flags holds the command-line flags for the `generate flatbuffers` command.
var (
	// flatbuffersDir specifies the input directory containing .fbs files.
	flatbuffersDir string
	// targetDir specifies the output directory for the generated code.
	targetDir string
	// lang specifies the target programming language.
	lang string
	// moduleOptions provides language-specific options, like a Go package path.
	moduleOptions string
)

// flatbuffersCmd represents the command to generate language-specific source code
// from FlatBuffer (.fbs) schema files.
var flatbuffersCmd = &cobra.Command{
	Use:   "flatbuffers",
	Short: "Generates language-specific source code from FlatBuffer schema files (.fbs)",
	Long: `The 'flatbuffers' command invokes the FlatBuffers compiler (flatc) to generate
source code for one or more target languages from your .fbs schema files.

You must specify the directory containing your .fbs files, the target language, and
the output directory for the generated code. For some languages, like Go, you may
also need to provide language-specific options.

Examples:
  # Generate Go code from FlatBuffer schemas
  buffman generate flatbuffers \
    --flatbuffers_dir=./path/to/fbs \
    --target_dir=./gen/go \
    --language=go \
    --module_options="github.com/your-org/your-project/gen/go"

  # Generate C++ code from FlatBuffer schemas
  buffman generate flatbuffers \
    --flatbuffers_dir=./schemas \
    --target_dir=./generated/cpp \
    --language=cpp`,
	Run: func(cmd *cobra.Command, args []string) {
		// The flags are already marked as required by cobra, but this provides a
		// fallback message for clarity.
		if lang == "" || targetDir == "" {
			fmt.Println("Error: please specify both --language and --target_dir flags.")
			cmd.Help()
			return
		}
		handleGenerate(lang, targetDir, flatbuffersDir, moduleOptions)
	},
}

// init registers and configures the flags for the `generate flatbuffers` command.
func init() {
	flatbuffersCmd.Flags().StringVarP(&flatbuffersDir, "flatbuffers_dir", "I", "", "Directory containing the source .fbs schema files")
	flatbuffersCmd.Flags().StringVarP(&targetDir, "target_dir", "o", "", "Output directory for the generated source code")
	flatbuffersCmd.Flags().StringVarP(&lang, "language", "l", "", "Target language for code generation (e.g., go, java, cpp, kotlin)")
	flatbuffersCmd.Flags().StringVarP(&moduleOptions, "module_options", "m", "", "Language-specific options (e.g., Go package path or Java package name)")

	// Mark flags as mandatory for command execution.
	flatbuffersCmd.MarkFlagRequired("language")
	flatbuffersCmd.MarkFlagRequired("target_dir")
}

// handleGenerate orchestrates the code generation process. It creates a generator
// instance, configures it with the provided options, and executes the code
// generation. It will print any errors and exit the program on failure.
func handleGenerate(lang, targetDir, flatbuffersDir, moduleOptions string) {
	g, err := generate.NewGenerate(generate.Flatbuffers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	opts := options.GenerateOptions{
		InputDir: flatbuffersDir,
		LanguageDetails: map[language.Language]options.LanguageGenerateOptions{
			language.Language(lang): {
				OutputDir: targetDir,
				Opts:      []string{parseModuleOptions(language.Language(lang), moduleOptions)},
			},
		},
	}

	if err := g.Generate(context.Background(), opts); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated %s files in %s\n", lang, targetDir)
}

// parseModuleOptions constructs a language-specific option string. It looks up the
// appropriate flag for a given language (e.g., `--go-pkg-prefix` for Go) and
// combines it with the user-provided option value. It returns an empty string
// if the option is empty or not defined for the language.
func parseModuleOptions(languageType language.Language, opt string) string {
	if opt == "" {
		return ""
	}

	commandOpts, ok := configuration.CommandOptionsMap[languageType]
	if !ok {
		return ""
	}
	// Example: returns "--go-pkg-prefix github.com/your-org/your-project/gen/go"
	return commandOpts["flatbuffer"].Options["package"].Flag + " " + opt
}
