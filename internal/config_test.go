package internal

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFile_Empty(t *testing.T) {
	reader := strings.NewReader("")

	config, err := LoadConfigFile(reader)
	assert.ErrorIs(t, err, io.EOF)
	assert.Nil(t, config)
}

func TestLoadConfigFile_EmptyProject(t *testing.T) {
	reader := strings.NewReader("project:\n")

	config, err := LoadConfigFile(reader)
	assert.NoError(t, err)

	if assert.NotNil(t, config) {
		if assert.NotNil(t, config.Project) {
			assert.Nil(t, config.Project.WorkDirectory)
			assert.Nil(t, config.Project.BuildDirectory)
			assert.Nil(t, config.Project.Environment)
		}
	}
}

func TestLoadConfigFile_Workflow_InvalidType(t *testing.T) {
	reader := strings.NewReader("project:\n  workflow: 3.15\n")

	config, err := LoadConfigFile(reader)
	assert.ErrorContains(t, err, "invalid workflow type: ")
	assert.Nil(t, config)
}

func TestLoadConfigFile_Workflow_InvalidStructure(t *testing.T) {
	reader := strings.NewReader("project:\n  workflow:\n    test: [1, 2]\n")

	config, err := LoadConfigFile(reader)
	assert.Error(t, err)
	assert.Nil(t, config)
}

//nolint:cyclop
func TestLoadConfigFile_Whole(t *testing.T) {
	reader := strings.NewReader(`
project:
    work-directory: /path/to/project
    build-directory: /path/to/build
    environment: docker:golang
    workflow:
        build: tool1
        configure: tool2
        test: tool3
        format: tool4
        lint: tool5
        run: tool6
        dependency: tool7
`,
	)

	config, err := LoadConfigFile(reader)
	assert.NoError(t, err)

	if assert.NotNil(t, config) {
		if assert.NotNil(t, config.Project) {
			if assert.NotNil(t, config.Project.WorkDirectory) {
				assert.Equal(t, "/path/to/project", *config.Project.WorkDirectory)
			}

			if assert.NotNil(t, config.Project.BuildDirectory) {
				assert.Equal(t, "/path/to/build", *config.Project.BuildDirectory)
			}

			if assert.NotNil(t, config.Project.Environment) {
				assert.Equal(t, "docker:golang", *config.Project.Environment)
			}

			if assert.NotNil(t, config.Project.Workflow) {
				workflow, ok := config.Project.Workflow.(*FileConfigProjectWorkflow)
				if assert.True(t, ok, "workflow should be of type FileConfigProjectWorkflow") {

					if assert.NotNil(t, workflow.Configure) {
						assert.Equal(t, "tool2", *workflow.Configure)
					}

					if assert.NotNil(t, workflow.Build) {
						assert.Equal(t, "tool1", *workflow.Build)
					}

					if assert.NotNil(t, workflow.Test) {
						assert.Equal(t, "tool3", *workflow.Test)
					}

					if assert.NotNil(t, workflow.Format) {
						assert.Equal(t, "tool4", *workflow.Format)
					}

					if assert.NotNil(t, workflow.Lint) {
						assert.Equal(t, "tool5", *workflow.Lint)
					}

					if assert.NotNil(t, workflow.Run) {
						assert.Equal(t, "tool6", *workflow.Run)
					}

					if assert.NotNil(t, workflow.Dependency) {
						assert.Equal(t, "tool7", *workflow.Dependency)
					}
				}
			}
		}
	}
}

//nolint:cyclop
func TestLoadConfigFile_Whole_WorkflowString(t *testing.T) {
	reader := strings.NewReader(`
project:
    work-directory: /path/to/project
    build-directory: /path/to/build
    environment: docker:golang
    workflow: go
`,
	)

	config, err := LoadConfigFile(reader)
	assert.NoError(t, err)

	if assert.NotNil(t, config) {
		if assert.NotNil(t, config.Project) {
			if assert.NotNil(t, config.Project.WorkDirectory) {
				assert.Equal(t, "/path/to/project", *config.Project.WorkDirectory)
			}

			if assert.NotNil(t, config.Project.BuildDirectory) {
				assert.Equal(t, "/path/to/build", *config.Project.BuildDirectory)
			}

			if assert.NotNil(t, config.Project.Environment) {
				assert.Equal(t, "docker:golang", *config.Project.Environment)
			}

			if assert.NotNil(t, config.Project.Workflow) {
				workflow, ok := config.Project.Workflow.(string)
				if assert.True(t, ok, "workflow should be of type string") {
					assert.Equal(t, "go", workflow)
				}
			}
		}
	}
}

func TestFileConfig_UpdateConfig_Empty(t *testing.T) {
	config := DefaultConfig()

	fileConfig := FileConfig{
		Project: FileConfigProject{},
	}

	fileConfig.UpdateConfig(&config)

	assert.Equal(t, DefaultConfig(), config)
}

//nolint:cyclop
func TestFileConfig_UpdateConfig(t *testing.T) {
	config := DefaultConfig()

	fileConfig := FileConfig{
		Project: FileConfigProject{
			WorkDirectory:  StrPtr("/project/work"),
			BuildDirectory: StrPtr("/project/build"),
			Environment:    StrPtr("env:arg"),
			Workflow: &FileConfigProjectWorkflow{
				Configure:  StrPtr("tool1"),
				Build:      StrPtr("tool2"),
				Test:       StrPtr("tool3"),
				Format:     StrPtr("tool4"),
				Lint:       StrPtr("tool5"),
				Run:        StrPtr("tool6"),
				Dependency: StrPtr("tool7"),
			},
		},
	}

	fileConfig.UpdateConfig(&config)

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.WorkDirectory) {
		assert.Equal(t, "/project/work", *config.WorkDirectory)
	}

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "/project/build", *config.BuildDirectory)
	}

	if assert.NotNil(t, config.Environment) {
		assert.Equal(t, "env:arg", *config.Environment)
	}

	if assert.NotNil(t, config.Workflow) {
		workflow, ok := config.Workflow.(*ConfigWorkflow)
		if assert.True(t, ok, "workflow should be of type ConfigWorkflow") {
			if assert.NotNil(t, workflow.Configure) {
				assert.Equal(t, "tool1", *workflow.Configure)
			}

			if assert.NotNil(t, workflow.Build) {
				assert.Equal(t, "tool2", *workflow.Build)
			}

			if assert.NotNil(t, workflow.Test) {
				assert.Equal(t, "tool3", *workflow.Test)
			}

			if assert.NotNil(t, workflow.Format) {
				assert.Equal(t, "tool4", *workflow.Format)
			}

			if assert.NotNil(t, workflow.Lint) {
				assert.Equal(t, "tool5", *workflow.Lint)
			}

			if assert.NotNil(t, workflow.Run) {
				assert.Equal(t, "tool6", *workflow.Run)
			}

			if assert.NotNil(t, workflow.Dependency) {
				assert.Equal(t, "tool7", *workflow.Dependency)
			}
		}
	}
}

func TestFileConfig_UpdateConfig_WorkflowString(t *testing.T) {
	config := DefaultConfig()

	fileConfig := FileConfig{
		Project: FileConfigProject{
			WorkDirectory:  StrPtr("/project/work"),
			BuildDirectory: StrPtr("/project/build"),
			Environment:    StrPtr("env:arg"),
			Workflow:       "my-workflow",
		},
	}

	fileConfig.UpdateConfig(&config)

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.WorkDirectory) {
		assert.Equal(t, "/project/work", *config.WorkDirectory)
	}

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "/project/build", *config.BuildDirectory)
	}

	if assert.NotNil(t, config.Environment) {
		assert.Equal(t, "env:arg", *config.Environment)
	}

	if assert.NotNil(t, config.Workflow) {
		workflow, ok := config.Workflow.(string)
		if assert.True(t, ok, "workflow should be of type string") {
			assert.Equal(t, "my-workflow", workflow)
		}
	}
}

func TestFileConfig_UpdateConfig_WorkflowInvalid(t *testing.T) {
	config := DefaultConfig()

	fileConfig := FileConfig{
		Project: FileConfigProject{
			WorkDirectory:  StrPtr("/project/work"),
			BuildDirectory: StrPtr("/project/build"),
			Environment:    StrPtr("env:arg"),
			Workflow:       42,
		},
	}

	assert.PanicsWithValue(t, "invalid workflow type: int", func() {
		fileConfig.UpdateConfig(&config)
	})
}
