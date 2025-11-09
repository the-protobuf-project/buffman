// package options defines configuration structures used across the Buffman CLI
// for code generation, parsing, and language-specific tasks. This package is
// separated from others to prevent circular import dependencies.
package options

import "github.com/machanirobotics/buffman/internal/generate/language"

// GenerateOptions holds the configuration for a code generation task.
// It specifies the input source directory and a map of language-specific
// generation details.
type GenerateOptions struct {
	// InputDir is the directory containing the source schema files.
	InputDir string
	// LanguageDetails maps each target language to its specific generation options.
	LanguageDetails map[language.Language]LanguageGenerateOptions
	// Verbose, if true, enables detailed logging during the generation process,
	// such as printing the commands being executed.
	Verbose bool
}

// ParseOptions holds the configuration for a parsing or conversion task.
// It specifies the input directory for source files and the output directory
// for the resulting artifacts.
type ParseOptions struct {
	// InputDir is the directory containing the source files to parse.
	InputDir string
	// OutputDir is the directory where the parsed output will be written.
	OutputDir string
	// Verbose, if true, enables detailed logging during the parsing process,
	// such as listing files as they are being converted.
	Verbose bool
}

// LanguageGenerateOptions holds language-specific settings for code generation.
// This includes the output directory and any additional command-line flags
// required by the external code generator.
type LanguageGenerateOptions struct {
	// OutputDir is the destination directory for the generated code for this language.
	OutputDir string
	// Opts is a slice of additional command-line options or flags for the generator.
	Opts []string
}
