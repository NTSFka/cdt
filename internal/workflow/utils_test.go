package workflow_test

import (
	"os"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"cdt/internal/workflow"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromConfig_Empty(t *testing.T) {
	config, err := workflow.FromConfig(internal.ConfigWorkflow{}, internal.Tools{})

	require.NoError(t, err)
	assert.Equal(t, &internal.Workflow{Name: "custom"}, config)
}

func TestFromConfig_Configure_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Configure: new("configure")}
	tools := internal.Tools{}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'configure' not found")
	assert.Nil(t, wf)
}

func TestFromConfig_Configure_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Configure: new("configure")}
	tools := internal.Tools{
		internal.NewExecutableTool("configure", "Test", "Test", internal.Tags{}, nil),
	}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'configure' doesn't support configuration")
	assert.Nil(t, wf)
}

func TestFromConfig_Configure(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectConfigurator
	}

	config := internal.ConfigWorkflow{Configure: new("configure")}
	testTool := &toolType{
		internal.MakeExecutableTool("configure", "Test", "Test", internal.Tags{}, nil),
		nil,
	}

	wf, err := workflow.FromConfig(config, internal.Tools{testTool})

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Equal(t, testTool, wf.Configurator)
	assert.Nil(t, wf.Builder)
	assert.Nil(t, wf.Tester)
	assert.Nil(t, wf.Formatter)
	assert.Nil(t, wf.Linter)
	assert.Nil(t, wf.Runner)
}

func TestFromConfig_Build_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Build: new("build")}
	tools := internal.Tools{}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'build' not found")
	assert.Nil(t, wf)
}

func TestFromConfig_Build_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Build: new("build")}
	tools := internal.Tools{
		internal.NewExecutableTool("build", "Test", "Test", internal.Tags{}, nil),
	}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'build' doesn't support building")
	assert.Nil(t, wf)
}

func TestFromConfig_Build(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectBuilder
	}

	config := internal.ConfigWorkflow{Build: new("build")}
	testTool := &toolType{
		internal.MakeExecutableTool("build", "Test", "Test", internal.Tags{}, nil),
		nil,
	}

	wf, err := workflow.FromConfig(config, internal.Tools{testTool})

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Nil(t, wf.Configurator)
	assert.Equal(t, testTool, wf.Builder)
	assert.Nil(t, wf.Tester)
	assert.Nil(t, wf.Formatter)
	assert.Nil(t, wf.Linter)
	assert.Nil(t, wf.Runner)
}

func TestFromConfig_Test_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Test: new("test")}
	tools := internal.Tools{}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'test' not found")
	assert.Nil(t, wf)
}

func TestFromConfig_Test_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Test: new("test")}
	tools := internal.Tools{
		internal.NewExecutableTool("test", "Test", "Test", internal.Tags{}, nil),
	}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'test' doesn't support testing")
	assert.Nil(t, wf)
}

func TestFromConfig_Test(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectTester
	}

	config := internal.ConfigWorkflow{Test: new("test")}
	testTool := &toolType{
		internal.MakeExecutableTool("test", "Test", "Test", internal.Tags{}, nil),
		nil,
	}

	wf, err := workflow.FromConfig(config, internal.Tools{testTool})

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Nil(t, wf.Configurator)
	assert.Nil(t, wf.Builder)
	assert.Equal(t, testTool, wf.Tester)
	assert.Nil(t, wf.Formatter)
	assert.Nil(t, wf.Linter)
	assert.Nil(t, wf.Runner)
}

func TestFromConfig_Format_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Format: new("format")}
	tools := internal.Tools{}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'format' not found")
	assert.Nil(t, wf)
}

func TestFromConfig_Format_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Format: new("format")}
	tools := internal.Tools{
		internal.NewExecutableTool("format", "Test", "Test", internal.Tags{}, nil),
	}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'format' doesn't support formatting")
	assert.Nil(t, wf)
}

func TestFromConfig_Format(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectFormatter
	}

	config := internal.ConfigWorkflow{Format: new("format")}
	testTool := &toolType{
		internal.MakeExecutableTool("format", "Test", "Test", internal.Tags{}, nil),
		nil,
	}

	wf, err := workflow.FromConfig(config, internal.Tools{testTool})

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Nil(t, wf.Configurator)
	assert.Nil(t, wf.Builder)
	assert.Nil(t, wf.Tester)
	assert.Equal(t, testTool, wf.Formatter)
	assert.Nil(t, wf.Linter)
	assert.Nil(t, wf.Runner)
}

func TestFromConfig_Lint_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Lint: new("lint")}
	tools := internal.Tools{}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'lint' not found")
	assert.Nil(t, wf)
}

func TestFromConfig_Lint_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Lint: new("lint")}
	tools := internal.Tools{
		internal.NewExecutableTool("lint", "Test", "Test", internal.Tags{}, nil),
	}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'lint' doesn't support linting")
	assert.Nil(t, wf)
}

func TestFromConfig_Lint(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectLinter
	}

	config := internal.ConfigWorkflow{Lint: new("lint")}
	testTool := &toolType{
		internal.MakeExecutableTool("lint", "Test", "Test", internal.Tags{}, nil),
		nil,
	}

	wf, err := workflow.FromConfig(config, internal.Tools{testTool})

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Nil(t, wf.Configurator)
	assert.Nil(t, wf.Builder)
	assert.Nil(t, wf.Tester)
	assert.Nil(t, wf.Formatter)
	assert.Equal(t, testTool, wf.Linter)
	assert.Nil(t, wf.Runner)
}

func TestFromConfig_Run_NotFound(t *testing.T) {
	config := internal.ConfigWorkflow{Run: new("run")}
	tools := internal.Tools{}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'run' not found")
	assert.Nil(t, wf)
}

func TestFromConfig_Run_NotSupported(t *testing.T) {
	config := internal.ConfigWorkflow{Run: new("run")}
	tools := internal.Tools{
		internal.NewExecutableTool("run", "Test", "Test", internal.Tags{}, nil),
	}

	wf, err := workflow.FromConfig(config, tools)

	require.EqualError(t, err, "tool 'run' doesn't support run")
	assert.Nil(t, wf)
}

func TestFromConfig_Run(t *testing.T) {
	type toolType struct {
		internal.ExecutableTool
		internal.ProjectRunner
	}

	config := internal.ConfigWorkflow{Run: new("run")}
	testTool := &toolType{
		internal.MakeExecutableTool("run", "Test", "Test", internal.Tags{}, nil),
		nil,
	}

	wf, err := workflow.FromConfig(config, internal.Tools{testTool})

	require.NoError(t, err)
	require.NotNil(t, wf)
	assert.Nil(t, wf.Configurator)
	assert.Nil(t, wf.Builder)
	assert.Nil(t, wf.Tester)
	assert.Nil(t, wf.Formatter)
	assert.Nil(t, wf.Linter)
	assert.Equal(t, testTool, wf.Runner)
}

func TestBuildProject_Default(t *testing.T) {
	project, err := workflow.CreateProject(internal.Config{RootDirectory: "dir1"}, internal.Tools{})

	require.NoError(t, err)

	require.NotNil(t, project)
	assert.Equal(t, "dir1", project.Info.Directory)
	assert.Nil(t, project.Workflow.Configurator)
	assert.Nil(t, project.Workflow.Builder)
	assert.Nil(t, project.Workflow.Tester)
	assert.Nil(t, project.Workflow.Formatter)
	assert.Nil(t, project.Workflow.Linter)
	assert.Nil(t, project.Workflow.Runner)
}

func TestBuildProject_Custom(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow: &internal.ConfigWorkflow{
			Configure: new("tool1"),
			Build:     new("tool2"),
			Test:      new("tool3"),
			Format:    new("tool4"),
			Lint:      new("tool5"),
			Run:       new("tool6"),
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

	project, err := workflow.CreateProject(
		config,
		internal.Tools{tool1, tool2, tool3, tool4, tool5, tool6},
	)

	require.NoError(t, err)

	require.NotNil(t, project)
	assert.Equal(t, "dir1", project.Info.Directory)
	assert.Equal(t, tool1, project.Workflow.Configurator)
	assert.Equal(t, tool2, project.Workflow.Builder)
	assert.Equal(t, tool3, project.Workflow.Tester)
	assert.Equal(t, tool4, project.Workflow.Formatter)
	assert.Equal(t, tool5, project.Workflow.Linter)
	assert.Equal(t, tool6, project.Workflow.Runner)
}

func TestBuildProject_Custom_Fail(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow: &internal.ConfigWorkflow{
			Configure: new("tool1"),
			Build:     new("tool2"),
			Test:      new("tool3"),
			Format:    new("tool4"),
			Lint:      new("tool5"),
			Run:       new("tool6"),
		},
	}
	project, err := workflow.CreateProject(config, internal.Tools{})

	require.EqualError(t, err, "tool 'tool1' not found")
	assert.Nil(t, project)
}

func TestBuildProject_Custom_UnknownWorkflow(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow:      99,
	}

	assert.PanicsWithValue(t, "unknown workflow type: int", func() {
		_, _ = workflow.CreateProject(config, internal.Tools{})
	})
}

func TestBuildProject_Custom_StringWorkflow_Unsupported(t *testing.T) {
	config := internal.Config{
		RootDirectory: "dir1",
		Workflow:      "my-workflow",
	}

	wf, err := workflow.CreateProject(config, internal.Tools{})

	require.EqualError(t, err, "workflow 'my-workflow' not found")
	assert.Nil(t, wf)
}

func TestBuildProject_Custom_StringWorkflow_Go(t *testing.T) {
	env := test.NewEnvironment(t)

	rootDir := t.TempDir()

	file, err := os.Create(filepath.Join(rootDir, "go.mod"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	config := internal.Config{RootDirectory: rootDir, Workflow: "go"}
	project, err := workflow.CreateProject(config, internal.Tools{
		tool.DetectGo(t.Context(), internal.DetectOptions{Environment: env}),
		tool.DetectNilAway(t.Context(), internal.DetectOptions{Environment: env}),
		tool.DetectGolangCILint(t.Context(), internal.DetectOptions{Environment: env}),
	})

	require.NoError(t, err)
	assert.NotNil(t, project)
}

func TestBuildProject_CMake(t *testing.T) {
	env := test.NewEnvironment(t)

	rootDir := t.TempDir()

	file, err := os.Create(filepath.Join(rootDir, "CMakeLists.txt"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	config := internal.Config{RootDirectory: rootDir}
	project, err := workflow.CreateProject(config, internal.Tools{
		tool.DetectCMake(t.Context(), internal.DetectOptions{Environment: env}),
		tool.DetectCTest(t.Context(), internal.DetectOptions{Environment: env}),
		tool.DetectClangFormat(t.Context(), internal.DetectOptions{Environment: env}),
		tool.DetectClangTidy(t.Context(), internal.DetectOptions{Environment: env}),
	})

	require.NoError(t, err)
	assert.NotNil(t, project)
}

func TestBuildProject_Go(t *testing.T) {
	env := test.NewEnvironment(t)

	rootDir := t.TempDir()

	file, err := os.Create(filepath.Join(rootDir, "go.mod"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	config := internal.Config{RootDirectory: rootDir}
	project, err := workflow.CreateProject(config, internal.Tools{
		tool.DetectGo(t.Context(), internal.DetectOptions{Environment: env}),
		tool.DetectNilAway(t.Context(), internal.DetectOptions{Environment: env}),
		tool.DetectGolangCILint(t.Context(), internal.DetectOptions{Environment: env}),
	})

	require.NoError(t, err)
	assert.NotNil(t, project)
}
