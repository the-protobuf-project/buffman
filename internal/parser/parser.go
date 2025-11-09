// Package parser provides interfaces and constructors for converting proto files
// into different target formats (e.g., FlatBuffers). It uses a factory and
// manager pattern to handle various parser implementations.
package parser

import (
	"context"
	"errors"

	"github.com/machanirobotics/buffman/internal/options"
	"github.com/machanirobotics/buffman/internal/parser/flatbuffers"
	"github.com/machanirobotics/buffman/internal/parser/nanobuffers"
)

// ParserType represents the type of parser to use (e.g., "flatbuffers").
type ParserType string

const (
	// Flatbuffers specifies the parser for converting from Protocol Buffers to FlatBuffers.
	Flatbuffers ParserType = "flatbuffers"

	Nanobuffers ParserType = "nanobuffers"
)

// Parser is the main interface for a file format conversion utility.
// Implementations are responsible for the logic of converting source files
// (like .proto) into a different schema format.
type Parser interface {
	// Parse executes the conversion process based on the provided options.
	// The context can be used to handle cancellation or timeouts.
	Parse(ctx context.Context, opts options.ParseOptions) error
}

// NewParser acts as a factory, returning a concrete implementation of the Parser
// interface based on the provided ParserType. It returns an error if the
// requested type is unsupported.
func NewParser(parserType ParserType) (Parser, error) {
	switch parserType {
	case Flatbuffers:
		return flatbuffers.NewFlatbuffersParser()
	case Nanobuffers:
		return nanobuffers.NewNanobuffersParser()
	default:
		return nil, errors.New("unsupported parser type")
	}
}

// Manager holds and manages a collection of registered Parser implementations.
type Manager struct {
	parsers map[ParserType]Parser
}

// NewManager initializes and returns a new, empty Manager.
func NewManager() *Manager {
	return &Manager{
		parsers: make(map[ParserType]Parser),
	}
}

// RegisterParsers creates and registers one or more parsers in the manager.
// It uses the NewParser factory and will return an error if any of the
// requested parser types are unsupported.
func (m *Manager) RegisterParsers(parserTypes ...ParserType) error {
	for _, p := range parserTypes {
		parser, err := NewParser(p)
		if err != nil {
			return err
		}
		m.parsers[p] = parser
	}
	return nil
}

// GetParser retrieves a registered parser by its type. It returns nil if the
// parser has not been registered.
func (m *Manager) GetParser(parserType ParserType) Parser {
	return m.parsers[parserType]
}

// ConvertAll executes the Parse method for all parsers that have corresponding
// options in the parserOpts map. It returns an error if a parser for a given
// option is not registered or if any parsing process fails.
func (m *Manager) ConvertAll(ctx context.Context, parserOpts map[ParserType]options.ParseOptions) error {
	for parserType, opts := range parserOpts {
		parser, exists := m.parsers[parserType]
		if !exists {
			return errors.New("unsupported parser")
		}
		if err := parser.Parse(ctx, opts); err != nil {
			return err
		}
	}
	return nil
}
