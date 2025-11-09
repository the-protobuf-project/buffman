package utilities

import (
	"os"
	"strings"
)

// InsertGeneratedComments adds the generated comment to the particular file
func InsertGeneratedComments(commentStr, langFilePath string) error {
	// Read the entire file content
	content, err := os.ReadFile(langFilePath)
	if err != nil {
		return &CommentInsertionError{Err: err}
	}

	// Convert to string and split into lines
	lines := strings.Split(string(content), "\n")

	// Remove the first line if there are any lines
	if len(lines) > 0 {
		lines = lines[1:]
	}

	// Join remaining lines back together
	remainingContent := strings.Join(lines, "\n")

	// Create new content: comment + remaining content (without first line)
	newContent := commentStr + remainingContent

	// Write the new content back to the file
	err = os.WriteFile(langFilePath, []byte(newContent), 0644)
	if err != nil {
		// remove the file if the commment generation fails
		if err := os.Remove(langFilePath); err != nil {
			return &CommentInsertionError{Err: err}
		}
		return err
	}

	return nil
}
