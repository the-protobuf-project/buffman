// Package configuration defines the configuration structures for the Buffman CLI.
// It uses the Koanf library to load, merge, and validate settings from a YAML
// configuration file. The configuration supports multiple buffer types (FlatBuffers,
// NanoBuffers, etc.) with multiple target languages, each having their own output
// directories and options.
package configuration

// Config is the top-level structure that maps directly to the buffman.yml file.
// It defines the version, input sources, and plugins for code generation.
type Config struct {
	Version string   `koanf:"version"` // Configuration file format version
	Inputs  []Input  `koanf:"inputs"`  // Source inputs for schema files
	Plugins []Plugin `koanf:"plugins"` // Code generation plugins configuration
}

// Input defines a source input which can be a local path or a remote git repository.
type Input struct {
	Name   string `koanf:"name"`   // Name of the input source
	Path   string `koanf:"path"`   // Local path to the directory containing schema files
	Remote string `koanf:"remote"` // Remote git repository URL (optional)
	Commit string `koanf:"commit"` // Commit hash for the remote repository (optional)
}

// Plugin represents a code generation plugin (e.g., flatbuffers, nanobuffers).
// Each plugin can generate code for multiple target languages with their own configurations.
type Plugin struct {
	Name      string             `koanf:"name"`      // Plugin name (e.g., "flatbuffers", "nanobuffers")
	Out       string             `koanf:"out"`       // Output directory of the schema
	Languages []LanguageGenerate `koanf:"languages"` // Target languages and their configurations
}

// LanguageGenerate defines the configuration for generating code in a specific language.
// It specifies the target language, output directory, and language-specific options.
type LanguageGenerate struct {
	Language string   `koanf:"language"` // Target language name (e.g., "go", "cpp", "java")
	Out      string   `koanf:"out"`      // Output directory for generated files
	Opt      []string `koanf:"opt"`      // Language-specific options as a list of strings
}

// GetPluginByName searches for a plugin with the specified name in the configuration.
// It returns a pointer to the plugin if found, or nil if no plugin with the given name exists.
func (c *Config) GetPluginByName(name string) *Plugin {
	for i, plugin := range c.Plugins {
		if plugin.Name == name {
			return &c.Plugins[i]
		}
	}
	return nil
}

// GetFlatbuffersPlugin returns the flatbuffers plugin configuration if it exists.
// This is a convenience method for accessing the most commonly used plugin.
// Returns nil if no flatbuffers plugin is configured.
func (c *Config) GetFlatbuffersPlugin() *Plugin {
	return c.GetPluginByName("flatbuffers")
}

// GetNanobuffersPlugin returns the nanobuffers plugin configuration if it exists.
// This is a convenience method for accessing the nanobuffers plugin.
// Returns nil if no nanobuffers plugin is configured.
func (c *Config) GetNanobuffersPlugin() *Plugin {
	return c.GetPluginByName("nanobuffers")
}

// HasAnyLanguage checks if the plugin has any languages configured.
// Returns true if at least one language is configured for this plugin, false otherwise.
func (p *Plugin) HasAnyLanguage() bool {
	return len(p.Languages) > 0
}

// GetLanguageConfig returns the configuration for a specific language within the plugin.
// It searches through the plugin's configured languages and returns the matching configuration.
// Returns nil if the specified language is not configured for this plugin.
func (p *Plugin) GetLanguageConfig(language string) *LanguageGenerate {
	for i, lang := range p.Languages {
		if lang.Language == language {
			return &p.Languages[i]
		}
	}
	return nil
}

// GetConfiguredLanguages returns a map of language names to their configurations.
// The map keys are language names (e.g., "go", "cpp") and values are their configurations.
// This method is useful for iterating over all configured languages for a plugin.
func (p *Plugin) GetConfiguredLanguages() map[string]*LanguageGenerate {
	languages := make(map[string]*LanguageGenerate)
	for i, lang := range p.Languages {
		languages[lang.Language] = &p.Languages[i]
	}
	return languages
}

// IsLocalInput returns true if the input is a local path (not a remote repository).
func (i *Input) IsLocalInput() bool {
	return i.Remote == ""
}

// IsRemoteInput returns true if the input is a remote git repository.
func (i *Input) IsRemoteInput() bool {
	return i.Remote != ""
}

// GetSourcePath returns the appropriate source path for the input.
// For local inputs, it returns the Path field. For remote inputs, it would
// typically return a path where the remote repository is cloned locally.
func (i *Input) GetSourcePath() string {
	if i.IsLocalInput() {
		return i.Path
	}
	// For remote inputs, this would typically be handled by a repository manager
	// that clones the repo to a temporary location
	return i.Path
}
