package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

// Config is an application configuration passed via flags or configuration file
type Config struct {
	// RootDirectory is a root directory of the project
	RootDirectory string

	// BuildDirectory is the project's build directory
	BuildDirectory *string

	// Environment defines an environment to use
	Environment *string
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		RootDirectory: ".",
	}
}

// FileConfig stores configuration from file
type FileConfig struct {
	Project FileConfigProject `yaml:"project"`
}

// UpdateConfig updates the given configuration by configuration from a file.
func (c *FileConfig) UpdateConfig(config *Config) {
	if c.Project.Directory != nil {
		config.RootDirectory = *c.Project.Directory
	}

	if c.Project.BuildDirectory != nil {
		config.BuildDirectory = c.Project.BuildDirectory
	}

	if c.Project.Environment != nil {
		config.Environment = c.Project.Environment
	}
}

// FileConfigProject stores configuration from a file: project part
type FileConfigProject struct {
	Directory      *string `yaml:"directory"`
	BuildDirectory *string `yaml:"build-directory"`
	Environment    *string `yaml:"environment"`
}

// LoadConfigFile loads configuration from a reader
func LoadConfigFile(reader io.Reader) (*FileConfig, error) {
	result := FileConfig{}

	err := yaml.NewDecoder(reader).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}

	return &result, nil
}
