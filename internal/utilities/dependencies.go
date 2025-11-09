package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetGoogleProtobufFiles clones a specific commit from the protobuf repository,
// extracts only the src/google/protobuf/ folder, and removes all non-.proto files.
func HandleGoogleProtobufFiles(googleProtbufDir string) error {

	// Extract only the src/google/protobuf/ folder and clean up
	targetFolder := "src/google/protobuf"
	if err := extractAndCleanProtobufFolder(googleProtbufDir, targetFolder); err != nil {
		return fmt.Errorf("failed to extract and clean protobuf folder: %w", err)
	}

	return nil
}

// extractAndCleanProtobufFolder extracts the target folder and removes everything else,
// then deletes all non-.proto files within the extracted folder.
func extractAndCleanProtobufFolder(repoDir, targetFolder string) error {
	sourcePath := filepath.Join(repoDir, targetFolder)

	// Check if the target folder exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("target folder %s does not exist in repository", targetFolder)
	}

	// Create a temporary directory to store the protobuf files
	tempDir := filepath.Join(repoDir, "temp_protobuf")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Copy only .proto files from the target folder to temp directory
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip subdirectories - only process files in the root target folder
		if info.IsDir() && path != sourcePath {
			return filepath.SkipDir
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".proto") {
			fileName := info.Name()

			// Read file contents to check for filtering criteria
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			contentStr := string(content)

			// Skip files with 'syntax = "proto2";' except 'descriptor.proto'
			if strings.Contains(contentStr, `syntax = "proto2";`) && !strings.Contains(fileName, "descriptor.proto") {
				return nil
			}

			if strings.Contains(fileName, "unittest") {
				return nil
			}

			// Skip files with 'edition' in the content
			if strings.Contains(contentStr, `edition = "202`) {
				return nil
			}

			// Calculate relative path from source
			relPath, err := filepath.Rel(sourcePath, path)
			if err != nil {
				return err
			}

			// Create destination path
			destPath := filepath.Join(tempDir, relPath)

			// Create destination directory if needed
			destDir := filepath.Dir(destPath)
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return err
			}

			// Copy the .proto file
			return copyFile(path, destPath)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy .proto files: %w", err)
	}

	// Remove everything in the repo directory except the temp folder
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return fmt.Errorf("failed to read repo directory: %w", err)
	}

	for _, entry := range entries {
		if entry.Name() != "temp_protobuf" {
			fullPath := filepath.Join(repoDir, entry.Name())
			if err := os.RemoveAll(fullPath); err != nil {
				return fmt.Errorf("failed to remove %s: %w", fullPath, err)
			}
		}
	}

	// Move contents from temp directory to root
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(tempDir, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(repoDir, relPath)
			destDir := filepath.Dir(destPath)

			if err := os.MkdirAll(destDir, 0755); err != nil {
				return err
			}

			return os.Rename(path, destPath)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to move files from temp directory: %w", err)
	}

	// Remove the temp directory
	return os.RemoveAll(tempDir)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}
