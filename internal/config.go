package internal

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// Config is an application configuration passed via flags or configuration file.
type Config struct {
	// RootDirectory is a root directory of the project
	RootDirectory string

	// OutputDirectory is the project's output directory
	OutputDirectory *string

	// Environment defines an environment to use
	Environment *string

	// Workflow defines tools to use in a workflow. Can be ConfigWorkflow or string.
	Workflow any

	// Tools stores custom tools configuration - tool name and executable path
	Tools ConfigTools
}

// DefaultConfig returns default configuration.
func DefaultConfig() Config {
	return Config{
		RootDirectory: ".",
	}
}

// ConfigTools stores custom tools configuration - tool name and executable path.
type ConfigTools map[string]string

// Get returns a tool executable path by its ID.
func (c ConfigTools) Get(id string, def string) string {
	if path, ok := c[id]; ok {
		return path
	} else {
		return def
	}
}

// ConfigWorkflow stores the project's workflow tools.
type ConfigWorkflow struct {
	Configure  *string
	Build      *string
	Test       *string
	Format     *string
	Lint       *string
	Run        *string
	Dependency *string
}

// FileConfig stores configuration from file.
type FileConfig struct {
	Tools   *FileConfigTools  `yaml:"tools"`
	Project FileConfigProject `yaml:"project"`
}

// UpdateConfig updates the given configuration by configuration from a file.
func (c *FileConfig) UpdateConfig(config *Config) {
	if c.Project.OutputDirectory != nil {
		config.OutputDirectory = c.Project.OutputDirectory
	}

	if c.Project.Environment != nil {
		config.Environment = c.Project.Environment
	}

	if c.Project.Workflow != nil {
		switch workflow := c.Project.Workflow.(type) {
		case *FileConfigProjectWorkflow:
			config.Workflow = &ConfigWorkflow{
				Configure:  workflow.Configure,
				Build:      workflow.Build,
				Test:       workflow.Test,
				Format:     workflow.Format,
				Lint:       workflow.Lint,
				Run:        workflow.Run,
				Dependency: workflow.Dependency,
			}
		case string:
			config.Workflow = workflow
		default:
			panic(fmt.Sprintf("invalid workflow type: %T", workflow))
		}
	}

	if c.Tools != nil {
		config.Tools = make(ConfigTools, len(*c.Tools))

		for k, v := range *c.Tools {
			config.Tools[k] = v
		}
	}
}

// FileConfigTools stores configuration from a file: tools part.
type FileConfigTools map[string]string

// FileConfigProject stores configuration from a file: project part.
type FileConfigProject struct {
	// OutputDirectory specifies a directory where output files can be store.
	OutputDirectory *string `yaml:"output-directory"`
	// Environment specifies which environment to run tools in for the given project
	Environment *string `yaml:"environment"`
	// Workflow specifies the project workflow, can be FileConfigProjectWorkflow or string
	Workflow any `yaml:"workflow"`
}

// FileConfigProjectWorkflow stores configuration from a file: project, workflow part.
type FileConfigProjectWorkflow struct {
	// Configure stores ID of a tool that be used as ProjectConfigurator
	Configure *string `yaml:"configure"`
	// Build stores ID of a tool that be used as ProjectBuilder
	Build *string `yaml:"build"`
	// Test stores ID of a tool that be used as ProjectTester
	Test *string `yaml:"test"`
	// Format stores ID of a tool that be used as ProjectFormatter
	Format *string `yaml:"format"`
	// Lint stores ID of a tool that be used as ProjectLinter
	Lint *string `yaml:"lint"`
	// Run stores ID of a tool that be used as ProjectRunner
	Run *string `yaml:"run"`
	// Run stores ID of a tool that be used as ProjectDependencyManager
	Dependency *string `yaml:"dependency"`
}

// LoadConfigFile loads configuration from a reader.
func LoadConfigFile(reader io.Reader) (*FileConfig, error) {
	result := FileConfig{}

	if err := yaml.NewDecoder(reader).Decode(&result); err != nil {
		return nil, fmt.Errorf("config load failed: %w", err)
	}

	// Check if Workflow is FileConfigProjectWorkflow or string
	if result.Project.Workflow != nil {
		switch value := result.Project.Workflow.(type) {
		case map[string]any:
			workflow := FileConfigProjectWorkflow{}

			var node yaml.Node

			// Can't fail
			_ = node.Encode(value)

			if err := node.Decode(&workflow); err != nil {
				return nil, err
			}

			result.Project.Workflow = &workflow
		case string:
			break
		default:
			return nil, fmt.Errorf("invalid workflow type: %T", value)
		}
	}

	return &result, nil
}
