package utilities

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ListFilesRelativeToRoot(rootPath, extenstion string) ([]string, error) {
	var paths []string
	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return &ListFilesError{Err: err}
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, extenstion) {
			relativePath := strings.TrimPrefix(path, rootPath)
			relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator))
			paths = append(paths, relativePath)
		}
		return nil
	})
	return paths, err
}

// getIncludePaths walks the entire directory tree starting from the generator's
// base flatbuffer directory. It collects all subdirectory paths and formats them
// as `-I <path>` strings. This allows the flatc compiler to resolve imports
// between .fbs files located in different subdirectories.
// It returns a slice of formatted include path strings or an error if the
// directory walk fails.
func GetIncludePaths(dir string, flag string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Add every directory encountered to the include paths.
		if d.IsDir() {
			paths = append(paths, flag+" "+path)
		}
		return nil
	})
	if err != nil {
		return nil, &IncludePathError{Err: err}
	}
	return paths, nil
}

// GenerateFiles writes files to disk based on the provided fileDetails map.
//
// fileDetails is a map where the key is the target file path and the value is
// the file content as a byte slice. The function creates any necessary
// directories for each file and writes the file content to disk.
//
// Returns an error if any directory or file cannot be created or written.
func GenerateFiles(fileDetails map[string][]byte) error {
	for path, content := range fileDetails {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
			return fmt.Errorf("failed to write file %s: %v", path, err)
		}
	}
	return nil
}

// ExcludeFiles filters out file paths that match any of the given patterns
func ExcludeFiles(filePaths []string, patterns ...string) []string {
	var result []string
	for _, filePath := range filePaths {
		shouldExclude := false
		for _, pattern := range patterns {
			// Check if the file path contains the pattern
			if strings.Contains(filePath, pattern) {
				shouldExclude = true
				break
			}
			// Also check if the pattern matches using filepath.Match for glob patterns
			if matched, err := filepath.Match(pattern, filePath); err == nil && matched {
				shouldExclude = true
				break
			}
			// Check if the pattern matches the base name of the file
			if matched, err := filepath.Match(pattern, filepath.Base(filePath)); err == nil && matched {
				shouldExclude = true
				break
			}
		}
		if !shouldExclude {
			result = append(result, filePath)
		}
	}
	return result
}
