package utilities

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// unzip extracts a zip archive to a destination directory.
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error closing zip reader: %v\n", err)
		}
	}()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip vulnerability: ensure file extraction is within the target directory.
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", path, err)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", path, err)
			}
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func(outFile *os.File) { // Ensure file is closed even on error
				if err := outFile.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Error closing extracted file: %v\n", err)
				}
			}(outFile)

			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer func() { // Ensure zip entry reader is closed
				if err := rc.Close(); err != nil {
					fmt.Fprintf(os.Stderr, "Error closing zip entry reader: %v\n", err)
				}
			}()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
