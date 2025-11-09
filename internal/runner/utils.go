package runner

import (
	"fmt"
	"path"
	"strings"

	"github.com/machanirobotics/buffman/internal/configuration"
	"github.com/machanirobotics/buffman/internal/generate"
	"github.com/machanirobotics/buffman/internal/generate/language"
	"github.com/machanirobotics/buffman/internal/options"
	"github.com/machanirobotics/buffman/internal/parser"
	"github.com/machanirobotics/buffman/internal/remote"
	"github.com/machanirobotics/buffman/internal/utilities"
)

// initializeRunner prepares the parser and generator managers and resolves
// any remote repositories defined in the configuration.
//
// It pulls remote inputs if specified, and registers the appropriate plugins
// and generators for the workflow.
func (r *runnerImpl) initializeRunner(config *configuration.Config) error {
	rem, err := remote.NewRemote(remote.Github)
	if err != nil {
		return err
	}
	source := r.getSource(config)

	for _, input := range config.Inputs {
		if input.Remote != "" {

			googleRepo := ""
			if strings.Contains(input.Name, "google") {
				googleRepo = "google"
			}

			remotePath := path.Join(source.Path, googleRepo, strings.Split(input.Remote, "/")[len(strings.Split(input.Remote, "/"))-1])
			if input.Name == source.Name {
				remotePath = input.Path
			}

			if err := rem.Pull(remote.PullOptions{Out: remotePath, Url: input.Remote, Commit: &input.Commit}); err != nil {
				return err
			}

			if strings.Contains(input.Remote, "https://github.com/protocolbuffers/protobuf") {
				if err := utilities.HandleGoogleProtobufFiles(remotePath); err != nil {
					return err
				}
			}
		}
	}

	parserManager := parser.NewManager()
	if err := parserManager.RegisterParsers(parser.Flatbuffers, parser.Nanobuffers); err != nil { // add other parsers once it is there, this is statically defined
		return err
	}

	generateManager := generate.NewManager()
	if err := generateManager.RegisterGenerate(generate.Flatbuffers, generate.Nanobuffers); err != nil {
		return err
	}

	r.Parser = parserManager
	r.Generate = generateManager
	r.ProtoDir = source.Path

	return nil
}

// getSource finds the primary source directory from the configuration. It iterates
// through the inputs and returns the one explicitly named "source". If no such
// input is found, it returns an empty Input struct.
func (r *runnerImpl) getSource(config *configuration.Config) configuration.Input {
	for _, input := range config.Inputs {
		if input.Name == "source" {
			return input
		}
	}
	return configuration.Input{}
}

// getParserOptions constructs a map of parser-specific options from the configuration.
// It iterates through the defined plugins, validates that each corresponds to a
// registered parser, and builds a ParseOptions object for it. It returns an
// error if any plugin name is not a supported parser type.
func (r *runnerImpl) getParserOptions(config *configuration.Config) (map[parser.ParserType]options.ParseOptions, error) {
	parserOptions := map[parser.ParserType]options.ParseOptions{}
	for _, plugin := range config.Plugins {
		parserType := parser.ParserType(plugin.Name)
		// Ensure that a parser for the given plugin name has been registered in the manager.
		if r.Parser.GetParser(parserType) == nil {
			return nil, fmt.Errorf("unsupported plugin name %s", plugin.Name)
		}
		parserOptions[parserType] = options.ParseOptions{
			InputDir:  r.ProtoDir,
			OutputDir: plugin.Out,
		}
	}
	return parserOptions, nil
}

// getGenerateOptions constructs a map of generator-specific options from the configuration.
// It iterates through the plugins, validates that each corresponds to a registered
// generator, and then delegates to getLangnguageOptions to build the detailed
// options for each language. It returns an error if a plugin is not a supported
// generator type or if language option parsing fails.
func (r *runnerImpl) getGenerateOptions(config *configuration.Config) (map[generate.GenerateType]options.GenerateOptions, error) {
	generateOptions := map[generate.GenerateType]options.GenerateOptions{}
	for _, plugin := range config.Plugins {
		generateType := generate.GenerateType(plugin.Name)
		// Ensure that a generator for the given plugin name has been registered in the manager.
		if r.Generate.GetGenerate(generateType) == nil {
			return nil, fmt.Errorf("unsupported invalid plugin name %s", plugin.Name)
		}
		languageOptions, err := r.getLangnguageOptions(plugin)
		if err != nil {
			return nil, err
		}
		// The input for generation is assumed to be the output directory of the conversion step.
		generateOptions[generateType] = options.GenerateOptions{
			InputDir:        plugin.Out,
			LanguageDetails: languageOptions,
		}
	}
	return generateOptions, nil
}

// getLangnguageOptions extracts and validates language-specific generation options
// from a single plugin configuration. It ensures each specified language is
// supported before creating its LanguageGenerateOptions. It returns an error if
// any language is unsupported.
func (r *runnerImpl) getLangnguageOptions(plugin configuration.Plugin) (map[language.Language]options.LanguageGenerateOptions, error) {
	languageOptions := map[language.Language]options.LanguageGenerateOptions{}
	for _, lang := range plugin.Languages {
		if !language.IsSupportedLanguage(lang.Language) {
			return nil, &language.UnsupportedLanguageError{}
		}
		langType := language.Language(lang.Language)
		languageOptions[langType] = options.LanguageGenerateOptions{
			OutputDir: lang.Out,
			Opts:      lang.Opt,
		}
	}
	return languageOptions, nil
}
