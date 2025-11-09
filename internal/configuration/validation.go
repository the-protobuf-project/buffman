package configuration

import (
	"fmt"
	"strings"
)

// Supported languages that can have options
var supportedLanguagesWithOptions = map[string]bool{
	"go":   true,
	"java": true,
}

// ValidateConfig serves as the main validation entry point, calling specific
// validation functions for each major section of the configuration.
// It ensures that all required fields are present and have valid values.
func ValidateConfig(c *Config) error {
	if err := validateVersion(c); err != nil {
		return err
	}

	if err := validateInputs(c); err != nil {
		return err
	}

	if err := validatePlugins(c); err != nil {
		return err
	}

	return nil
}

// validateVersion checks that the configuration version is present and valid.
func validateVersion(c *Config) error {
	if c.Version == "" {
		return NewInvalidConfigError("version", "version is required")
	}
	return nil
}

// validateInputs checks that at least one input is configured and all inputs are valid.
// It also enforces that there is a compulsory input named "source" with a non-empty path.
func validateInputs(c *Config) error {
	if len(c.Inputs) == 0 {
		return NewInvalidConfigError("inputs", "at least one input must be configured")
	}

	foundSource := false

	for i, input := range c.Inputs {
		if input.Name == "source" {
			foundSource = true
			if input.Path == "" && input.Remote == "" {
				return NewInvalidConfigErrorWithIndex("inputs", "path or remote is required for input named 'source'", i)
			}
			// For "source", remote and commit are optional, so no error if missing
		}

		if err := validateInput(&input, i); err != nil {
			return err
		}
	}

	if !foundSource {
		return NewInvalidConfigError("inputs", "an input named 'source' must be configured")
	}

	return nil
}

// validateInput checks that an individual input configuration is valid.
func validateInput(input *Input, index int) error {
	if input.Name == "source" {
		// For the compulsory "source" input, path must be present, remote and commit optional
		if input.Path == "" {
			return NewInvalidConfigErrorWithIndex("inputs", "path is required for input named 'source'", index)
		}
		return nil
	}

	// For other inputs, validate as usual
	if input.Name == "" {
		return NewInvalidConfigErrorWithIndex("inputs", "name is required", index)
	}

	// Check if it's a local input
	if input.IsLocalInput() {
		if input.Path == "" {
			return NewInvalidConfigErrorWithIndex("inputs",
				fmt.Sprintf("path is required for local input '%s'", input.Name), index)
		}
		if input.Commit != "" {
			return NewInvalidConfigErrorWithIndex("inputs",
				fmt.Sprintf("commit cannot be specified for local input '%s'", input.Name), index)
		}
	} else {
		// It's a remote input
		if input.Remote == "" {
			return NewInvalidConfigErrorWithIndex("inputs",
				fmt.Sprintf("remote URL is required for remote input '%s'", input.Name), index)
		}
		if input.Commit == "" {
			return NewInvalidConfigErrorWithIndex("inputs",
				fmt.Sprintf("commit is required for remote input '%s'", input.Name), index)
		}
	}

	return nil
}

// validatePlugins checks that at least one plugin is configured and all plugins are valid.
func validatePlugins(c *Config) error {
	if len(c.Plugins) == 0 {
		return NewInvalidConfigError("plugins", "at least one plugin must be configured")
	}

	for i, plugin := range c.Plugins {
		if err := validatePlugin(&plugin, i); err != nil {
			return err
		}
	}

	return nil
}

// validatePlugin checks that a plugin configuration is valid.
// It ensures that the plugin has a name and at least one language configured,
// and that each language configuration is complete and valid.
func validatePlugin(plugin *Plugin, index int) error {
	if plugin.Name == "" {
		return NewInvalidConfigErrorWithIndex("plugins", "name is required", index)
	}

	if plugin.Out == "" {
		return NewInvalidConfigErrorWithIndex("plugins",
			fmt.Sprintf("output path is required for plugin '%s'", plugin.Name), index)
	}

	if !plugin.HasAnyLanguage() && plugin.Name != "nanobuffers" {
		return NewInvalidConfigErrorWithIndex("plugins",
			fmt.Sprintf("plugin '%s' must have at least one language configured", plugin.Name), index)
	}

	for j, lang := range plugin.Languages {
		if err := validateLanguage(&lang, plugin.Name, index, j); err != nil {
			return err
		}
	}

	return nil
}

// validateLanguage checks that a language configuration is valid.
// It validates the language name, output directory, and language-specific options.
func validateLanguage(lang *LanguageGenerate, pluginName string, pluginIndex, langIndex int) error {
	if lang.Language == "" {
		return NewInvalidConfigErrorWithIndex("plugins",
			fmt.Sprintf("language is required for plugin '%s'", pluginName), pluginIndex)
	}

	if lang.Out == "" {
		return NewInvalidConfigErrorWithIndex("plugins",
			fmt.Sprintf("output directory is required for language '%s' in plugin '%s'",
				lang.Language, pluginName), pluginIndex)
	}

	// Validate and transform language-specific options
	if err := validateLanguageOptions(lang, pluginName); err != nil {
		return err
	}

	return nil
}

// validateLanguageOptions checks that options are only provided for supported languages.
// It also validates and transforms required Go and Java options.
func validateLanguageOptions(lang *LanguageGenerate, pluginName string) error {
	hasOptions := len(lang.Opt) > 0
	supportsOptions := supportedLanguagesWithOptions[lang.Language]

	if hasOptions && !supportsOptions {
		return NewUnsupportedLanguageOptionsError(lang.Language, pluginName, lang.Opt)
	}

	// Add the new validation and transformation
	if err := validateAndTransformLanguageOptions(lang, pluginName); err != nil {
		return err
	}

	return nil
}

// validateAndTransformLanguageOptions enforces and transforms Go and Java options.
func validateAndTransformLanguageOptions(lang *LanguageGenerate, pluginName string) error {
	switch lang.Language {
	case "go":
		found := false
		for i, opt := range lang.Opt {
			if strings.HasPrefix(opt, "go_package=") {
				parts := strings.SplitN(opt, "=", 2)
				if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
					return fmt.Errorf("invalid go_package option in plugin '%s': %q", pluginName, opt)
				}
				lang.Opt[i] = "go-module-name " + strings.TrimSpace(parts[1])
				found = true
			}
		}
		if !found {
			return fmt.Errorf("go_package option is required for language 'go' in plugin '%s'", pluginName)
		}
	case "java":
		for i, opt := range lang.Opt {
			if strings.HasPrefix(opt, "java_package_prefix=") {
				parts := strings.SplitN(opt, "=", 2)
				if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
					return fmt.Errorf("invalid java_package_prefix option in plugin '%s': %q", pluginName, opt)
				}
				lang.Opt[i] = "java-package-prefix " + strings.TrimSpace(parts[1])
			}
		}
	}
	return nil
}
