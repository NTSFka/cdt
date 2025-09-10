package project

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildProject_Default(t *testing.T) {
	project, err := BuildProject(internal.Config{RootDirectory: "dir1"}, internal.Tools{})

	assert.NoError(t, err)

	if assert.NotNil(t, project) {
		assert.Equal(t, "dir1", project.Info.Directory)
		assert.Nil(t, project.Workflow.Configurator)
		assert.Nil(t, project.Workflow.Builder)
		assert.Nil(t, project.Workflow.Tester)
		assert.Nil(t, project.Workflow.Formatter)
		assert.Nil(t, project.Workflow.Linter)
		assert.Nil(t, project.Workflow.Runner)
	}
}

func TestBuildProject_Custom(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow: &internal.ConfigWorkflow{
			Configure: internal.StrPtr("tool1"),
			Build:     internal.StrPtr("tool2"),
			Test:      internal.StrPtr("tool3"),
			Format:    internal.StrPtr("tool4"),
			Lint:      internal.StrPtr("tool5"),
			Run:       internal.StrPtr("tool6"),
		},
	}
	tool1 := &struct {
		internal.ExecutableTool
		internal.ProjectConfigurator
	}{internal.MakeExecutableTool("tool1", "Test", "Test", internal.Tags{}, nil), nil}
	tool2 := &struct {
		internal.ExecutableTool
		internal.ProjectBuilder
	}{internal.MakeExecutableTool("tool2", "Test", "Test", internal.Tags{}, nil), nil}
	tool3 := &struct {
		internal.ExecutableTool
		internal.ProjectTester
	}{internal.MakeExecutableTool("tool3", "Test", "Test", internal.Tags{}, nil), nil}
	tool4 := &struct {
		internal.ExecutableTool
		internal.ProjectFormatter
	}{internal.MakeExecutableTool("tool4", "Test", "Test", internal.Tags{}, nil), nil}
	tool5 := &struct {
		internal.ExecutableTool
		internal.ProjectLinter
	}{internal.MakeExecutableTool("tool5", "Test", "Test", internal.Tags{}, nil), nil}
	tool6 := &struct {
		internal.ExecutableTool
		internal.ProjectRunner
	}{internal.MakeExecutableTool("tool6", "Test", "Test", internal.Tags{}, nil), nil}

	project, err := BuildProject(config, internal.Tools{tool1, tool2, tool3, tool4, tool5, tool6})

	assert.NoError(t, err)

	if assert.NotNil(t, project) {
		assert.Equal(t, "dir1", project.Info.Directory)
		assert.Equal(t, tool1, project.Workflow.Configurator)
		assert.Equal(t, tool2, project.Workflow.Builder)
		assert.Equal(t, tool3, project.Workflow.Tester)
		assert.Equal(t, tool4, project.Workflow.Formatter)
		assert.Equal(t, tool5, project.Workflow.Linter)
		assert.Equal(t, tool6, project.Workflow.Runner)
	}
}

func TestBuildProject_Custom_Fail(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow: &internal.ConfigWorkflow{
			Configure: internal.StrPtr("tool1"),
			Build:     internal.StrPtr("tool2"),
			Test:      internal.StrPtr("tool3"),
			Format:    internal.StrPtr("tool4"),
			Lint:      internal.StrPtr("tool5"),
			Run:       internal.StrPtr("tool6"),
		},
	}
	project, err := BuildProject(config, internal.Tools{})

	assert.EqualError(t, err, "tool 'tool1' not found")
	assert.Nil(t, project)
}

func TestBuildProject_Custom_UnknownWorkflow(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow:      99,
	}

	assert.PanicsWithValue(t, "unknown workflow type: int", func() {
		_, _ = BuildProject(config, internal.Tools{})
	})
}

func TestBuildProject_Custom_StringWorkflow_Unsupported(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow:      "my-workflow",
	}

	workflow, err := BuildProject(config, internal.Tools{})

	assert.EqualError(t, err, "workflow 'my-workflow' not found")
	assert.Nil(t, workflow)
}

func TestBuildProject_Custom_StringWorkflow_Go(t *testing.T) {
	env := test.NewEnvironment(t)

	rootDir := t.TempDir()

	if _, err := os.Create(filepath.Join(rootDir, "go.mod")); err != nil {
		assert.NoError(t, err)
	}

	config := internal.Config{RootDirectory: rootDir, Workflow: "go"}
	project, err := BuildProject(config, internal.Tools{
		tool.DetectGo(env),
		tool.DetectGolangCILint(env),
	})

	assert.NoError(t, err)
	assert.NotNil(t, project)
}

func TestBuildProject_CMake(t *testing.T) {
	env := test.NewEnvironment(t)

	rootDir := t.TempDir()

	if _, err := os.Create(filepath.Join(rootDir, "CMakeLists.txt")); err != nil {
		assert.NoError(t, err)
	}

	config := internal.Config{RootDirectory: rootDir}
	project, err := BuildProject(config, internal.Tools{
		tool.DetectCMake(env),
		tool.DetectCTest(env),
		tool.DetectClangFormat(env),
		tool.DetectClangTidy(env),
	})

	assert.NoError(t, err)
	assert.NotNil(t, project)
}

func TestBuildProject_Go(t *testing.T) {
	env := test.NewEnvironment(t)

	rootDir := t.TempDir()

	if _, err := os.Create(filepath.Join(rootDir, "go.mod")); err != nil {
		assert.NoError(t, err)
	}

	config := internal.Config{RootDirectory: rootDir}
	project, err := BuildProject(config, internal.Tools{
		tool.DetectGo(env),
		tool.DetectGolangCILint(env),
	})

	assert.NoError(t, err)
	assert.NotNil(t, project)
}
