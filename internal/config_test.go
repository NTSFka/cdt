package internal_test

import (
	"cdt/internal"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFile_Empty(t *testing.T) {
	reader := strings.NewReader("")

	config, err := internal.LoadConfigFile(reader)
	require.ErrorIs(t, err, io.EOF)
	assert.Nil(t, config)
}

func TestLoadConfigFile_EmptyProject(t *testing.T) {
	reader := strings.NewReader("project:\n")

	config, err := internal.LoadConfigFile(reader)
	require.NoError(t, err)

	require.NotNil(t, config)
	require.NotNil(t, config.Project)
	assert.Nil(t, config.Project.WorkDirectory)
	assert.Nil(t, config.Project.BuildDirectory)
	assert.Nil(t, config.Project.Environment)
}

func TestLoadConfigFile_Workflow_InvalidType(t *testing.T) {
	reader := strings.NewReader("project:\n  workflow: 3.15\n")

	config, err := internal.LoadConfigFile(reader)
	require.ErrorContains(t, err, "invalid workflow type: ")
	assert.Nil(t, config)
}

func TestLoadConfigFile_Workflow_InvalidStructure(t *testing.T) {
	reader := strings.NewReader("project:\n  workflow:\n    test: [1, 2]\n")

	config, err := internal.LoadConfigFile(reader)
	require.Error(t, err)
	require.Nil(t, config)
}

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

	config, err := internal.LoadConfigFile(reader)
	require.NoError(t, err)

	require.NotNil(t, config)
	require.NotNil(t, config.Project)
	require.NotNil(t, config.Project.WorkDirectory)
	assert.Equal(t, "/path/to/project", *config.Project.WorkDirectory)

	require.NotNil(t, config.Project.BuildDirectory)
	assert.Equal(t, "/path/to/build", *config.Project.BuildDirectory)

	require.NotNil(t, config.Project.Environment)
	assert.Equal(t, "docker:golang", *config.Project.Environment)

	require.NotNil(t, config.Project.Workflow)
	workflow, ok := config.Project.Workflow.(*internal.FileConfigProjectWorkflow)
	require.True(t, ok, "workflow should be of type FileConfigProjectWorkflow")

	require.NotNil(t, workflow.Configure)
	assert.Equal(t, "tool2", *workflow.Configure)

	require.NotNil(t, workflow.Build)
	assert.Equal(t, "tool1", *workflow.Build)

	require.NotNil(t, workflow.Test)
	assert.Equal(t, "tool3", *workflow.Test)

	require.NotNil(t, workflow.Format)
	assert.Equal(t, "tool4", *workflow.Format)

	require.NotNil(t, workflow.Lint)
	assert.Equal(t, "tool5", *workflow.Lint)

	require.NotNil(t, workflow.Run)
	assert.Equal(t, "tool6", *workflow.Run)

	require.NotNil(t, workflow.Dependency)
	assert.Equal(t, "tool7", *workflow.Dependency)
}

func TestLoadConfigFile_Whole_WorkflowString(t *testing.T) {
	reader := strings.NewReader(`
project:
    work-directory: /path/to/project
    build-directory: /path/to/build
    environment: docker:golang
    workflow: go
`,
	)

	config, err := internal.LoadConfigFile(reader)
	require.NoError(t, err)
	require.NotNil(t, config)
	require.NotNil(t, config.Project)

	require.NotNil(t, config.Project.WorkDirectory)
	assert.Equal(t, "/path/to/project", *config.Project.WorkDirectory)

	require.NotNil(t, config.Project.BuildDirectory)
	assert.Equal(t, "/path/to/build", *config.Project.BuildDirectory)

	require.NotNil(t, config.Project.Environment)
	assert.Equal(t, "docker:golang", *config.Project.Environment)

	require.NotNil(t, config.Project.Workflow)
	workflow, ok := config.Project.Workflow.(string)
	require.True(t, ok, "workflow should be of type string")
	assert.Equal(t, "go", workflow)
}

func TestFileConfig_UpdateConfig_Empty(t *testing.T) {
	config := internal.DefaultConfig()

	fileConfig := internal.FileConfig{
		Project: internal.FileConfigProject{},
	}

	fileConfig.UpdateConfig(&config)

	assert.Equal(t, internal.DefaultConfig(), config)
}

func TestFileConfig_UpdateConfig(t *testing.T) {
	config := internal.DefaultConfig()

	fileConfig := internal.FileConfig{
		Project: internal.FileConfigProject{
			WorkDirectory:  internal.StrPtr("/project/work"),
			BuildDirectory: internal.StrPtr("/project/build"),
			Environment:    internal.StrPtr("env:arg"),
			Workflow: &internal.FileConfigProjectWorkflow{
				Configure:  internal.StrPtr("tool1"),
				Build:      internal.StrPtr("tool2"),
				Test:       internal.StrPtr("tool3"),
				Format:     internal.StrPtr("tool4"),
				Lint:       internal.StrPtr("tool5"),
				Run:        internal.StrPtr("tool6"),
				Dependency: internal.StrPtr("tool7"),
			},
		},
	}

	fileConfig.UpdateConfig(&config)

	assert.Equal(t, ".", config.RootDirectory)

	require.NotNil(t, config.WorkDirectory)
	assert.Equal(t, "/project/work", *config.WorkDirectory)

	require.NotNil(t, config.BuildDirectory)
	assert.Equal(t, "/project/build", *config.BuildDirectory)

	require.NotNil(t, config.Environment)
	assert.Equal(t, "env:arg", *config.Environment)

	require.NotNil(t, config.Workflow)
	workflow, ok := config.Workflow.(*internal.ConfigWorkflow)
	require.True(t, ok, "workflow should be of type ConfigWorkflow")
	require.NotNil(t, workflow.Configure)
	assert.Equal(t, "tool1", *workflow.Configure)

	require.NotNil(t, workflow.Build)
	assert.Equal(t, "tool2", *workflow.Build)

	require.NotNil(t, workflow.Test)
	assert.Equal(t, "tool3", *workflow.Test)

	require.NotNil(t, workflow.Format)
	assert.Equal(t, "tool4", *workflow.Format)

	require.NotNil(t, workflow.Lint)
	assert.Equal(t, "tool5", *workflow.Lint)

	require.NotNil(t, workflow.Run)
	assert.Equal(t, "tool6", *workflow.Run)

	require.NotNil(t, workflow.Dependency)
	assert.Equal(t, "tool7", *workflow.Dependency)
}

func TestFileConfig_UpdateConfig_WorkflowString(t *testing.T) {
	config := internal.DefaultConfig()

	fileConfig := internal.FileConfig{
		Project: internal.FileConfigProject{
			WorkDirectory:  internal.StrPtr("/project/work"),
			BuildDirectory: internal.StrPtr("/project/build"),
			Environment:    internal.StrPtr("env:arg"),
			Workflow:       "my-workflow",
		},
	}

	fileConfig.UpdateConfig(&config)

	assert.Equal(t, ".", config.RootDirectory)

	require.NotNil(t, config.WorkDirectory)
	assert.Equal(t, "/project/work", *config.WorkDirectory)

	require.NotNil(t, config.BuildDirectory)
	assert.Equal(t, "/project/build", *config.BuildDirectory)

	require.NotNil(t, config.Environment)
	assert.Equal(t, "env:arg", *config.Environment)

	require.NotNil(t, config.Workflow)
	workflow, ok := config.Workflow.(string)
	require.True(t, ok, "workflow should be of type string")
	assert.Equal(t, "my-workflow", workflow)
}

func TestFileConfig_UpdateConfig_WorkflowInvalid(t *testing.T) {
	config := internal.DefaultConfig()

	fileConfig := internal.FileConfig{
		Project: internal.FileConfigProject{
			WorkDirectory:  internal.StrPtr("/project/work"),
			BuildDirectory: internal.StrPtr("/project/build"),
			Environment:    internal.StrPtr("env:arg"),
			Workflow:       42,
		},
	}

	assert.PanicsWithValue(t, "invalid workflow type: int", func() {
		fileConfig.UpdateConfig(&config)
	})
}
