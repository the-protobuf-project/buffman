// package generate provides a factory and manager for different code generators,
// allowing for a unified interface to generate source code from various
// schema types.
package generate

import (
	"context"
	"errors"

	"github.com/machanirobotics/buffman/internal/generate/flatbuffers"
	"github.com/machanirobotics/buffman/internal/generate/nanobuffers"
	"github.com/machanirobotics/buffman/internal/options"
)

// GenerateType defines the type of code generator, used by the factory to
// create a specific instance.
type GenerateType string

const (
	// Flatbuffers represents the FlatBuffers code generator.
	Flatbuffers GenerateType = "flatbuffers"

	// Nanobuffers represents the NanoBuffers code generator
	Nanobuffers GenerateType = "nanobuffers"
	// Additional generator types can be added here in the future.
)

// Generate defines the interface for a code generator. Each implementation is
// responsible for transforming schema files into language-specific source code.
type Generate interface {
	// Generate executes the code generation process based on the provided options.
	// The context can be used to manage cancellation or timeouts.
	Generate(ctx context.Context, opts options.GenerateOptions) error
}

// NewGenerate acts as a factory for creating Generate instances. It takes a
// GenerateType and returns the corresponding generator implementation. It returns an
// error if the requested type is unsupported.
func NewGenerate(generateType GenerateType) (Generate, error) {
	switch generateType {
	case Flatbuffers:
		return flatbuffers.NewFlatbuffersGenerate(), nil

	case Nanobuffers:
		return nanobuffers.NewNanoBuffersGenerate(), nil
	default:
		return nil, errors.New("unsupported generate type")
	}
}

// Manager holds and manages a collection of registered code generators.
type Manager struct {
	generators map[GenerateType]Generate
}

// NewManager initializes and returns an empty Manager.
func NewManager() *Manager {
	return &Manager{generators: make(map[GenerateType]Generate)}
}

// RegisterGenerate creates and registers one or more generators in the manager.
// It uses the NewGenerate factory and will return an error if any of the
// requested generator types are unsupported.
func (m *Manager) RegisterGenerate(generateTypes ...GenerateType) error {
	for _, g := range generateTypes {
		generate, err := NewGenerate(g)
		if err != nil {
			return err
		}
		m.generators[g] = generate
	}
	return nil
}

// GetGenerate retrieves a registered generator by its type. It returns nil if the
// generator has not been registered.
func (m *Manager) GetGenerate(generateType GenerateType) Generate {
	return m.generators[generateType]
}

// GenerateAll executes the Generate method for all generators that have
// corresponding options in the generateOpts map. It returns an error if a
// generator for a given option is not registered or if any generation process
// fails.
func (m *Manager) GenerateAll(ctx context.Context, generateOpts map[GenerateType]options.GenerateOptions) error {
	for generateType, opts := range generateOpts {
		generate, exists := m.generators[generateType]
		if !exists {
			return errors.New("unsupported generate type")
		}
		if err := generate.Generate(ctx, opts); err != nil {
			return err
		}
	}
	return nil
}
