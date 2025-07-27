package project

import (
	"cdt/internal"
	"github.com/stretchr/testify/assert"
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
