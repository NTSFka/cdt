package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

// ConfigEnvironment stores configuration for environment
type ConfigEnvironment struct {
	ToolName string
	Argument string
}

// Config is an application configuration passed via flags or configuration file
type Config struct {
	// RootDirectory is a root directory of the project
	RootDirectory string

	// BuildDirectory is the project's build directory
	BuildDirectory *string

	// Environment defines an environment to use
	Environment *ConfigEnvironment
}

// FileConfig stores configuration from file
type FileConfig struct {
	Project fileConfigProject `yaml:"project"`
}

type fileConfigProject struct {
	Directory      *string `yaml:"directory"`
	BuildDirectory *string `yaml:"build-directory"`
	Environment    *string `yaml:"environment"`
}

// LoadConfigFile loads configuration from file
func LoadConfigFile(reader io.Reader) (*FileConfig, error) {
	result := FileConfig{}

	err := yaml.NewDecoder(reader).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}

	return &result, nil
}
