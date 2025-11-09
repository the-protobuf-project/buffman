package nanobuffers

import (
	"context"

	"github.com/machanirobotics/buffman/internal/options"
)

type NanobuffersGenerate struct {
}

func NewNanoBuffersGenerate() *NanobuffersGenerate {
	return &NanobuffersGenerate{}
}

// a dummy to preven additional loc during execution
func (n *NanobuffersGenerate) Generate(ctx context.Context, opts options.GenerateOptions) error {
	return nil
}
