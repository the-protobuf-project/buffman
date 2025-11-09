// this file removes the google api protos

package flatbuffers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/machanirobotics/buffman/internal/utilities"
)

var excludePatterns = []string{"google/api", "google/protobuf/descriptor.proto"}

// clearGoogleAPI removes google API imports from proto files and writes cleaned files to the cleanedDir.
//
// Returns an error if proto files cannot be listed, cleaned, or written.
func (c *FlatbuffersParser) clearGoogleAPI(ctx context.Context) error {
	// List all proto files in the source directory
	protoFiles, err := utilities.ListFilesRelativeToRoot(c.protoDir, ".proto")
	if err != nil {
		return err
	}

	if len(protoFiles) == 0 {
		return errors.New("no proto files found")
	}

	// Remove google/api files and descriptor.proto as flatc does not support its conversion
	filteredFiles := utilities.ExcludeFiles(protoFiles, excludePatterns...)

	// Process files to remove Google API imports
	fileDetails, err := utilities.RemoveGoogleAPI(ctx, c.compiler, filteredFiles, c.cleanedDir)
	if err != nil {
		return err
	}

	// Create files and necessary directories
	return utilities.GenerateFiles(fileDetails)
}

// convertProtoFile converts a single proto file to a FlatBuffers schema using flatc.
//
// It creates the necessary output directory structure and post-processes the generated
// .fbs file to fix include statements.
//
// Returns an error if conversion or post-processing fails.
func (c *FlatbuffersParser) convertProtoFile(ctx context.Context, protoFile string) error {
	// Calculate relative path for the proto file
	sourceFilePath := path.Join(c.cleanedDir, protoFile)
	relPath, err := filepath.Rel(c.cleanedDir, sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	// Create target directory structure
	targetDir := filepath.Join(c.flatbufferDir, filepath.Dir(relPath))
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %s: %w", targetDir, err)
	}

	// Execute flatc command to convert proto to FlatBuffers
	if err := c.executeFlatcCommand(ctx, sourceFilePath, targetDir); err != nil {
		// Clean up on failure
		if removeErr := os.RemoveAll(targetDir); removeErr != nil {
			return removeErr
		}
		return err
	}

	// Post-process the generated .fbs file
	return c.postProcessFBSFile(protoFile, targetDir)
}

// executeFlatcCommand runs the flatc command to convert proto to FlatBuffers
func (c *FlatbuffersParser) executeFlatcCommand(ctx context.Context, sourceFile, targetDir string) error {
	cmd := exec.CommandContext(ctx, "flatc",
		"--proto",
		"-I", c.cleanedDir,
		"-o", targetDir,
		sourceFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("flatc command failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// postProcessFBSFile fixes include statements in the generated .fbs file
func (c *FlatbuffersParser) postProcessFBSFile(protoFile, targetDir string) error {
	fbsFileName := strings.TrimSuffix(filepath.Base(protoFile), ".proto") + ".fbs"
	fbsFilePath := filepath.Join(targetDir, fbsFileName)

	if err := c.fixFBSIncludes(fbsFilePath); err != nil {
		fmt.Printf("  ⚠️  Warning: failed to fix includes in %s: %v\n", fbsFileName, err)
	}

	return nil
}

// fixFBSIncludes modifies include statements in a generated FlatBuffers schema file.
// It replaces `import "file.proto";` with `include "file.fbs";`.
func (c *FlatbuffersParser) fixFBSIncludes(fbsFile string) error {
	content, err := os.ReadFile(fbsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("FBS file not found: %s", fbsFile)
		}
		return fmt.Errorf("failed to read FBS file: %w", err)
	}

	originalContent := string(content)
	importPattern := regexp.MustCompile(`import\s+"([^"]+)\.proto"\s*;`)
	modifiedContent := importPattern.ReplaceAllString(originalContent, `include "$1.fbs";`)

	if modifiedContent != originalContent {
		if err := os.WriteFile(fbsFile, []byte(modifiedContent), 0644); err != nil {
			return fmt.Errorf("failed to write modified FBS file: %w", err)
		}
	}
	return nil
}
