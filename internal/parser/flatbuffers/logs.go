package flatbuffers

import "fmt"

// FileProcessingReport holds the results and logs of a file conversion process.
type FileProcessingReport struct {
	LogMessages  []string
	FailedFiles  []string
	SuccessCount int
	ErrorCount   int
}

// NewFileProcessingReport initializes a new, empty report.
func NewFileProcessingReport() *FileProcessingReport {
	return &FileProcessingReport{
		LogMessages: make([]string, 0),
		FailedFiles: make([]string, 0),
	}
}

// recordSuccess updates the report for a successful file conversion.
func (r *FileProcessingReport) recordSuccess(message string) {
	r.LogMessages = append(r.LogMessages, message)
	r.SuccessCount++
}

// recordFailure updates the report for a failed file conversion.
func (r *FileProcessingReport) recordFailure(file, message string) {
	r.LogMessages = append(r.LogMessages, message)
	r.FailedFiles = append(r.FailedFiles, file)
	r.ErrorCount++
}

// BuildFullLog assembles the initial logs, per-file logs, and a final summary.
func (r *FileProcessingReport) BuildFullLog(initialLogs []string, flatbufferDir string) []string {
	summaryLogs := []string{
		"\n📊Conversion Summary:",
		fmt.Sprintf("   Successful: %d files", r.SuccessCount),
		fmt.Sprintf("   Failed: %d files", r.ErrorCount),
		fmt.Sprintf("   Output directory: %s", flatbufferDir),
	}
	fullLog := append(initialLogs, r.LogMessages...)
	return append(fullLog, summaryLogs...)
}

// logVerbose prints log messages and a list of failed files if verbose is true.
func (c *FlatbuffersParser) logVerbose(verbose bool, messages, failedFiles []string) {
	if !verbose {
		return
	}
	for _, msg := range messages {
		fmt.Println(msg)
	}
	if len(failedFiles) > 0 {
		fmt.Println("\n❌ Failed Files:")
		for _, file := range failedFiles {
			fmt.Printf("   - %s\n", file)
		}
	}
}
