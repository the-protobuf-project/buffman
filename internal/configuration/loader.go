package configuration

import (
	"fmt"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// LoadConfig reads the specified YAML configuration file, applies defaults, validates
// the settings, and unmarshals them into a Config struct.
//
// The process follows these steps:
// 1. Resolve the absolute path of the configuration file.
// 2. Load the hard-coded default values.
// 3. Load the user-provided YAML file, which overrides the defaults.
// 4. Unmarshal the final configuration into the Config struct.
// 5. Validate the struct to ensure all required fields are present.
//
// It returns a populated Config struct or an error if any step fails.
func LoadConfig(filename string) (*Config, error) {
	k := koanf.New(".")

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, NewConfigurationLoadError(filename,
			fmt.Errorf("failed to resolve file path: %w", err))
	}

	if err := loadDefaults(k); err != nil {
		return nil, NewConfigurationLoadError(filename,
			fmt.Errorf("failed to load default configuration: %w", err))
	}

	if err := k.Load(file.Provider(absPath), yaml.Parser()); err != nil {
		return nil, NewConfigurationLoadError(filename, err)
	}

	var config Config
	if err := k.Unmarshal("", &config); err != nil {
		return nil, NewConfigurationLoadError(filename,
			fmt.Errorf("failed to unmarshal config: %w", err))
	}

	if err := ValidateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// loadDefaults sets up the default configuration values using a confmap provider.
// These defaults ensure that the application has a sane starting point and users
// only need to specify values they wish to override.
// The defaults include a basic input directory and an empty plugins array.
func loadDefaults(k *koanf.Koanf) error {
	defaults := map[string]any{
		"version": "v1",
		"inputs": []map[string]any{
			{
				"name": "default",
				"path": "protobuf",
			},
		},
		"plugins": []map[string]any{},
	}

	return k.Load(confmap.Provider(defaults, "."), nil)
}
