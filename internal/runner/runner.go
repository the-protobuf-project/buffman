// Package runner defines the high-level interface for executing tasks based on a config file.
// It abstracts the core logic of the Buffman application, orchestrating the
// conversion and generation workflows for multiple buffer types and languages.
package runner

import (
	"context"
	"fmt"

	"github.com/machanirobotics/buffman/internal/configuration"
	"github.com/machanirobotics/buffman/internal/generate"
	"github.com/machanirobotics/buffman/internal/parser"
)

// Runner defines the high-level interface for executing tasks based on a config file.
// It abstracts the core logic of the Buffman application, orchestrating the
// conversion and generation workflows across multiple plugins and languages.
type Runner interface {
	// Run loads the configuration from the given file path and executes all defined
	// plugin tasks sequentially. The context can be used to cancel the entire run.
	// It returns an error if any step fails.
	Run(ctx context.Context, filePath string) error
}

// NewRunner creates a new instance of the default Runner implementation.
// The returned Runner can process configurations with multiple plugins
// and generate code for various target languages.
func NewRunner() Runner {
	return &runnerImpl{}
}

// runnerImpl is the concrete implementation of the Runner interface.
// It handles the orchestration of conversion and generation tasks
// based on the loaded configuration.
type runnerImpl struct {
	ProtoDir string
	Parser   *parser.Manager
	Generate *generate.Manager
}

// Run loads the configuration from the specified file path and executes all
// configured plugins in sequence. It processes each input directory and
// runs all configured plugins for code generation.
//
// Run returns an error if any configuration loading, parsing, or generation step fails.
func (r *runnerImpl) Run(ctx context.Context, filePath string) error {
	config, err := configuration.LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	if err := r.initializeRunner(config); err != nil {
		return err
	}

	parseOptions, err := r.getParserOptions(config)
	if err != nil {
		return err
	}
	if err := r.Parser.ConvertAll(ctx, parseOptions); err != nil {
		return err
	}

	generateOptions, err := r.getGenerateOptions(config)
	if err != nil {
		return err
	}

	if err := r.Generate.GenerateAll(ctx, generateOptions); err != nil {
		return err
	}

	return nil
}
