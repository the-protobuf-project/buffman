// package flatbuffers implements the code generator for FlatBuffer (.fbs) schemas.
package flatbuffers

import (
	"context"
	"strings"

	"github.com/machanirobotics/buffman/internal/generate/language"
	"github.com/machanirobotics/buffman/internal/options"
)

// FlatbufferGenerate is the concrete implementation of the generate.Generate
// interface for FlatBuffers.
type FlatbufferGenerate struct {
}

// NewFlatbuffersGenerate creates and returns a new instance of FlatbufferGenerate.
func NewFlatbuffersGenerate() *FlatbufferGenerate {
	return &FlatbufferGenerate{}
}

// Generate executes the FlatBuffers code generation process (flatc). It iterates
// over the target languages defined in the options, constructs the necessary
// command-line arguments, and invokes the generation logic for each language.
//
// It returns an error if language metadata cannot be found or if the underlying
// generation command fails.
func (f *FlatbufferGenerate) Generate(ctx context.Context, opts options.GenerateOptions) error {
	for lang, langDetails := range opts.LanguageDetails {
		metadata, err := language.GetMetadata(lang)
		if err != nil {
			return err
		}

		// The generateLanguageFile function (defined elsewhere in this package)
		// is responsible for executing the actual flatc command.
		err = generateLanguageFile(
			opts.InputDir,
			langDetails.OutputDir,
			strings.Join(langDetails.Opts, " "),
			metadata,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
