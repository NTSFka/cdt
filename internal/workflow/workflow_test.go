package workflow

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFromConfig_Empty(t *testing.T) {
	config, err := FromConfig(internal.ConfigWorkflow{}, internal.Tools{})

	assert.NoError(t, err)
	assert.Equal(t, &internal.Workflow{}, config)
}

func TestFromConfig_Configure_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Configure: internal.StrPtr("configure")}
	tools := internal.Tools{}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'configure' not found")
	assert.Nil(t, tool)
}

func TestFromConfig_Configure_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Configure: internal.StrPtr("configure")}
	tools := internal.Tools{
		internal.NewExecutableTool("configure", "Test", "Test", nil),
	}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'configure' doesn't support configuration")
	assert.Nil(t, tool)
}

func TestFromConfig_Configure(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectConfigurator
	}

	config := internal.ConfigWorkflow{Configure: internal.StrPtr("configure")}
	tool := &toolType{
		internal.MakeExecutableTool("configure", "Test", "Test", nil),
		nil,
	}

	workflow, err := FromConfig(config, internal.Tools{tool})

	assert.NoError(t, err)
	if assert.NotNil(t, workflow) {
		assert.Equal(t, tool, workflow.Configurator)
		assert.Nil(t, workflow.Builder)
		assert.Nil(t, workflow.Tester)
		assert.Nil(t, workflow.Formatter)
		assert.Nil(t, workflow.Linter)
		assert.Nil(t, workflow.Runner)
	}
}

func TestFromConfig_Build_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Build: internal.StrPtr("build")}
	tools := internal.Tools{}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'build' not found")
	assert.Nil(t, tool)
}

func TestFromConfig_Build_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Build: internal.StrPtr("build")}
	tools := internal.Tools{
		internal.NewExecutableTool("build", "Test", "Test", nil),
	}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'build' doesn't support building")
	assert.Nil(t, tool)
}

func TestFromConfig_Build(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectBuilder
	}

	config := internal.ConfigWorkflow{Build: internal.StrPtr("build")}
	tool := &toolType{
		internal.MakeExecutableTool("build", "Test", "Test", nil),
		nil,
	}

	workflow, err := FromConfig(config, internal.Tools{tool})

	assert.NoError(t, err)
	if assert.NotNil(t, workflow) {
		assert.Nil(t, workflow.Configurator)
		assert.Equal(t, tool, workflow.Builder)
		assert.Nil(t, workflow.Tester)
		assert.Nil(t, workflow.Formatter)
		assert.Nil(t, workflow.Linter)
		assert.Nil(t, workflow.Runner)
	}
}

func TestFromConfig_Test_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Test: internal.StrPtr("test")}
	tools := internal.Tools{}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'test' not found")
	assert.Nil(t, tool)
}

func TestFromConfig_Test_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Test: internal.StrPtr("test")}
	tools := internal.Tools{
		internal.NewExecutableTool("test", "Test", "Test", nil),
	}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'test' doesn't support testing")
	assert.Nil(t, tool)
}

func TestFromConfig_Test(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectTester
	}

	config := internal.ConfigWorkflow{Test: internal.StrPtr("test")}
	tool := &toolType{
		internal.MakeExecutableTool("test", "Test", "Test", nil),
		nil,
	}

	workflow, err := FromConfig(config, internal.Tools{tool})

	assert.NoError(t, err)
	if assert.NotNil(t, workflow) {
		assert.Nil(t, workflow.Configurator)
		assert.Nil(t, workflow.Builder)
		assert.Equal(t, tool, workflow.Tester)
		assert.Nil(t, workflow.Formatter)
		assert.Nil(t, workflow.Linter)
		assert.Nil(t, workflow.Runner)
	}
}

func TestFromConfig_Format_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Format: internal.StrPtr("format")}
	tools := internal.Tools{}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'format' not found")
	assert.Nil(t, tool)
}

func TestFromConfig_Format_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Format: internal.StrPtr("format")}
	tools := internal.Tools{
		internal.NewExecutableTool("format", "Test", "Test", nil),
	}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'format' doesn't support formatting")
	assert.Nil(t, tool)
}

func TestFromConfig_Format(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectFormatter
	}

	config := internal.ConfigWorkflow{Format: internal.StrPtr("format")}
	tool := &toolType{
		internal.MakeExecutableTool("format", "Test", "Test", nil),
		nil,
	}

	workflow, err := FromConfig(config, internal.Tools{tool})

	assert.NoError(t, err)
	if assert.NotNil(t, workflow) {
		assert.Nil(t, workflow.Configurator)
		assert.Nil(t, workflow.Builder)
		assert.Nil(t, workflow.Tester)
		assert.Equal(t, tool, workflow.Formatter)
		assert.Nil(t, workflow.Linter)
		assert.Nil(t, workflow.Runner)
	}
}

func TestFromConfig_Lint_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Lint: internal.StrPtr("lint")}
	tools := internal.Tools{}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'lint' not found")
	assert.Nil(t, tool)
}

func TestFromConfig_Lint_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Lint: internal.StrPtr("lint")}
	tools := internal.Tools{
		internal.NewExecutableTool("lint", "Test", "Test", nil),
	}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'lint' doesn't support linting")
	assert.Nil(t, tool)
}

func TestFromConfig_Lint(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectLinter
	}

	config := internal.ConfigWorkflow{Lint: internal.StrPtr("lint")}
	tool := &toolType{
		internal.MakeExecutableTool("lint", "Test", "Test", nil),
		nil,
	}

	workflow, err := FromConfig(config, internal.Tools{tool})

	assert.NoError(t, err)
	if assert.NotNil(t, workflow) {
		assert.Nil(t, workflow.Configurator)
		assert.Nil(t, workflow.Builder)
		assert.Nil(t, workflow.Tester)
		assert.Nil(t, workflow.Formatter)
		assert.Equal(t, tool, workflow.Linter)
		assert.Nil(t, workflow.Runner)
	}
}

func TestFromConfig_Run_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Run: internal.StrPtr("run")}
	tools := internal.Tools{}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'run' not found")
	assert.Nil(t, tool)
}

func TestFromConfig_Run_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Run: internal.StrPtr("run")}
	tools := internal.Tools{
		internal.NewExecutableTool("run", "Test", "Test", nil),
	}

	tool, err := FromConfig(config, tools)

	assert.EqualError(t, err, "tool 'run' doesn't support run")
	assert.Nil(t, tool)
}

func TestFromConfig_Run(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectRunner
	}

	config := internal.ConfigWorkflow{Run: internal.StrPtr("run")}
	tool := &toolType{
		internal.MakeExecutableTool("run", "Test", "Test", nil),
		nil,
	}

	workflow, err := FromConfig(config, internal.Tools{tool})

	assert.NoError(t, err)
	if assert.NotNil(t, workflow) {
		assert.Nil(t, workflow.Configurator)
		assert.Nil(t, workflow.Builder)
		assert.Nil(t, workflow.Tester)
		assert.Nil(t, workflow.Formatter)
		assert.Nil(t, workflow.Linter)
		assert.Equal(t, tool, workflow.Runner)
	}
}

func TestBuildProject_Default(t *testing.T) {
	project, err := BuildProject(internal.Config{RootDirectory: "dir1"}, internal.Tools{})

	assert.NoError(t, err)

	if assert.NotNil(t, project) {
		assert.Equal(t, "dir1", project.RootDirectory())
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
	}{internal.MakeExecutableTool("tool1", "Test", "Test", nil), nil}
	tool2 := &struct {
		internal.ExecutableTool
		internal.ProjectBuilder
	}{internal.MakeExecutableTool("tool2", "Test", "Test", nil), nil}
	tool3 := &struct {
		internal.ExecutableTool
		internal.ProjectTester
	}{internal.MakeExecutableTool("tool3", "Test", "Test", nil), nil}
	tool4 := &struct {
		internal.ExecutableTool
		internal.ProjectFormatter
	}{internal.MakeExecutableTool("tool4", "Test", "Test", nil), nil}
	tool5 := &struct {
		internal.ExecutableTool
		internal.ProjectLinter
	}{internal.MakeExecutableTool("tool5", "Test", "Test", nil), nil}
	tool6 := &struct {
		internal.ExecutableTool
		internal.ProjectRunner
	}{internal.MakeExecutableTool("tool6", "Test", "Test", nil), nil}

	project, err := BuildProject(config, internal.Tools{tool1, tool2, tool3, tool4, tool5, tool6})

	assert.NoError(t, err)

	if assert.NotNil(t, project) {
		assert.Equal(t, "dir1", project.RootDirectory())
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
		tool.DetectClangFormat(env, nil),
		tool.DetectClangTidy(env, nil),
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
