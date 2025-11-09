package nanobuffers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/machanirobotics/buffman/internal/options"
	"github.com/machanirobotics/buffman/internal/utilities"
)

const nanobuffersCmd = "nanopb %s -D %s %s" // include paths, directory, glob path of the protofiles

type NanobuffersParser struct {
}

func NewNanobuffersParser() (*NanobuffersParser, error) {
	return &NanobuffersParser{}, nil
}

func (n *NanobuffersParser) Parse(ctx context.Context, opts options.ParseOptions) error {
	includePaths, err := utilities.GetIncludePaths(opts.InputDir, "-I")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return err
	}
	if err := n.executeCommand(opts, strings.Join(includePaths, " ")); err != nil {
		return err
	}
	return nil
}

func (n *NanobuffersParser) executeCommand(opts options.ParseOptions, includePaths string) error {
	cmdStr := fmt.Sprintf("shopt -s globstar; %s",
		fmt.Sprintf(nanobuffersCmd, includePaths, opts.OutputDir, path.Join(opts.InputDir, "**/*.proto")))

	cmd := exec.Command("bash", "-c", cmdStr)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution error: %v, ouput %s", err, output)
	}

	return nil

}
