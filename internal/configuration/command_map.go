// Package configuration provides a static mapping between abstract concepts,
// like "package," and the concrete command-line flags required by external
// tools like the FlatBuffers compiler (flatc). This allows the application's
// logic to remain decoupled from tool-specific implementations.
package configuration

import (
	"github.com/machanirobotics/buffman/internal"
	"github.com/machanirobotics/buffman/internal/generate/language"
)

// Subcommand represents an abstract command-line option concept, such as
// defining a package.
type Subcommand string

const (
	// Pkg represents the concept of defining a package or module.
	Pkg Subcommand = "package"
)

// ToolOption defines a specific command-line flag and its description for
// an external tool.
type ToolOption struct {
	// Flag is the command-line flag (e.g., "--go-module-name").
	Flag string
	// Description explains the purpose and usage of the flag.
	Description string
}

// CommandOptions aggregates the configuration for a given tool, including its
// version and a map of its supported options.
type CommandOptions struct {
	// Version specifies the version of the tool these options apply to.
	Version string
	// Options maps abstract concepts (Subcommand) to concrete tool flags (ToolOption).
	Options map[Subcommand]ToolOption
}

// CommandOptionsMap provides a static, multi-level mapping to translate abstract
// configuration concepts into concrete command-line flags for external tools.
//
// The map is structured as follows:
//   - Level 1 Key: language.Language (e.g., Go, Java)
//   - Level 2 Key: string (a tool or generator name, e.g., "flatbuffer")
//
// This allows the application to look up tool-specific flags programmatically.
// For example, to find the flag for defining a Go package for FlatBuffers, the
// lookup would be: CommandOptionsMap[language.Go]["flatbuffer"].Options[Pkg].
var CommandOptionsMap = map[language.Language]map[string]CommandOptions{
	language.Go: {
		"buffman": {
			Version: internal.Version,
			Options: map[Subcommand]ToolOption{
				Pkg: {
					Description: `Specifies the Go package name, often used as a file option. ` +
						`Example: go_package = "github.com/machanirobotics/buffman/examples/go/fb"`,
				},
			},
		},
		"flatbuffer": {
			// This corresponds to the version of the flatc compiler.
			Version: "2.0.0",
			Options: map[Subcommand]ToolOption{
				Pkg: {
					Flag: "--go-module-name",
					Description: `Specifies the Go module name when generating FlatBuffers code. ` +
						`Example: --go-module-name github.com/machanirobotics/buffman/examples/go/fb`,
				},
			},
		},
		// Future tools like nanobuffers can be added here.
	},
	language.Java: {
		"flatbuffer": {
			Version: "2.0.0",
			Options: map[Subcommand]ToolOption{
				Pkg: {
					Flag: "--java-package-prefix",
					Description: `Specifies the package prefix for generated Java code. ` +
						`Example: --java-package-prefix com.machanirobotics.fb`,
				},
			},
		},
	},
	// Additional languages can be added here.
}
