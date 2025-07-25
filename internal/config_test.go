package internal

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
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
				if assert.NotNil(t, config.Project.Workflow.Configure) {
					assert.Equal(t, "tool2", *config.Project.Workflow.Configure)
				}

				if assert.NotNil(t, config.Project.Workflow.Build) {
					assert.Equal(t, "tool1", *config.Project.Workflow.Build)
				}

				if assert.NotNil(t, config.Project.Workflow.Test) {
					assert.Equal(t, "tool3", *config.Project.Workflow.Test)
				}

				if assert.NotNil(t, config.Project.Workflow.Format) {
					assert.Equal(t, "tool4", *config.Project.Workflow.Format)
				}

				if assert.NotNil(t, config.Project.Workflow.Lint) {
					assert.Equal(t, "tool5", *config.Project.Workflow.Lint)
				}

				if assert.NotNil(t, config.Project.Workflow.Run) {
					assert.Equal(t, "tool6", *config.Project.Workflow.Run)
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

func TestFileConfig_UpdateConfig(t *testing.T) {
	strPtr := func(s string) *string {
		return &s
	}

	config := DefaultConfig()

	fileConfig := FileConfig{
		Project: FileConfigProject{
			WorkDirectory:  strPtr("/project/work"),
			BuildDirectory: strPtr("/project/build"),
			Environment:    strPtr("env:arg"),
			Workflow: &FileConfigProjectWorkflow{
				Configure: strPtr("tool1"),
				Build:     strPtr("tool2"),
				Test:      strPtr("tool3"),
				Format:    strPtr("tool4"),
				Lint:      strPtr("tool5"),
				Run:       strPtr("tool6"),
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
		if assert.NotNil(t, config.Workflow.Configure) {
			assert.Equal(t, "tool1", *config.Workflow.Configure)
		}

		if assert.NotNil(t, config.Workflow.Build) {
			assert.Equal(t, "tool2", *config.Workflow.Build)
		}

		if assert.NotNil(t, config.Workflow.Test) {
			assert.Equal(t, "tool3", *config.Workflow.Test)
		}

		if assert.NotNil(t, config.Workflow.Format) {
			assert.Equal(t, "tool4", *config.Workflow.Format)
		}

		if assert.NotNil(t, config.Workflow.Lint) {
			assert.Equal(t, "tool5", *config.Workflow.Lint)
		}

		if assert.NotNil(t, config.Workflow.Run) {
			assert.Equal(t, "tool6", *config.Workflow.Run)
		}
	}
}
