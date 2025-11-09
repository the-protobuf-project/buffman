package configuration

import "fmt"

// InvalidConfigurationError represents an error in the configuration file.
// It provides detailed information about what went wrong during configuration
// validation or loading.
type InvalidConfigurationError struct {
	Field   string // The configuration field that caused the error
	Message string // Detailed error message
	Index   int    // Index in array/slice if applicable (-1 if not applicable)
}

// Error implements the error interface for InvalidConfigurationError.
func (e *InvalidConfigurationError) Error() string {
	if e.Index >= 0 {
		return fmt.Sprintf("invalid configuration at %s[%d]: %s", e.Field, e.Index, e.Message)
	}
	return fmt.Sprintf("invalid configuration for %s: %s", e.Field, e.Message)
}

// NewInvalidConfigError creates a new InvalidConfigurationError for a specific field.
func NewInvalidConfigError(field, message string) *InvalidConfigurationError {
	return &InvalidConfigurationError{
		Field:   field,
		Message: message,
		Index:   -1,
	}
}

// NewInvalidConfigErrorWithIndex creates a new InvalidConfigurationError for an array/slice element.
func NewInvalidConfigErrorWithIndex(field, message string, index int) *InvalidConfigurationError {
	return &InvalidConfigurationError{
		Field:   field,
		Message: message,
		Index:   index,
	}
}

// UnsupportedLanguageOptionsError represents an error when a language has
// options configured but the language doesn't support options.
type UnsupportedLanguageOptionsError struct {
	Language string   // The language that doesn't support options
	Plugin   string   // The plugin where this error occurred
	Options  []string // The options that were provided
}

// Error implements the error interface for UnsupportedLanguageOptionsError.
func (e *UnsupportedLanguageOptionsError) Error() string {
	return fmt.Sprintf("language '%s' in plugin '%s' does not support options, but %d option(s) were provided: %v",
		e.Language, e.Plugin, len(e.Options), e.Options)
}

// NewUnsupportedLanguageOptionsError creates a new UnsupportedLanguageOptionsError.
func NewUnsupportedLanguageOptionsError(language, plugin string, options []string) *UnsupportedLanguageOptionsError {
	return &UnsupportedLanguageOptionsError{
		Language: language,
		Plugin:   plugin,
		Options:  options,
	}
}

// ConfigurationLoadError represents an error that occurred while loading the configuration file.
type ConfigurationLoadError struct {
	Filename string // The configuration file that failed to load
	Cause    error  // The underlying error that caused the failure
}

// Error implements the error interface for ConfigurationLoadError.
func (e *ConfigurationLoadError) Error() string {
	return fmt.Sprintf("failed to load configuration file '%s': %v", e.Filename, e.Cause)
}

// Unwrap returns the underlying error for error unwrapping.
func (e *ConfigurationLoadError) Unwrap() error {
	return e.Cause
}

// NewConfigurationLoadError creates a new ConfigurationLoadError.
func NewConfigurationLoadError(filename string, cause error) *ConfigurationLoadError {
	return &ConfigurationLoadError{
		Filename: filename,
		Cause:    cause,
	}
}
