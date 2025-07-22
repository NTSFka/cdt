package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runConfigure(configurator internal.ProjectConfigurator, args ...string) error {
	return test.RunCommand(NewConfigureCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Configurator: configurator,
			},
		},
	}, args...)
}

func runConfigureTool(configurator internal.Tool, args ...string) error {
	return test.RunCommand(NewConfigureCommand(), internal.Context{
		Project: internal.Project{},
		Tools: []internal.Tool{
			configurator,
		},
	}, args...)
}

func TestConfigure_NotSupported(t *testing.T) {
	err := runConfigure(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support configuration", err.Error())
	}
}

func TestConfigure_Configure_Success(t *testing.T) {
	configurator := test.ProjectConfigurator{}
	configurator.On("Configure", mock.Anything, []string{}).Return(nil)

	err := runConfigure(&configurator)

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestConfigure_Configure_Failure(t *testing.T) {
	configurator := test.ProjectConfigurator{}
	configurator.On("Configure", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runConfigure(&configurator)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestConfigure_Tool_Success(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectConfigurator
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectConfigurator{},
	}
	linter.On("Configure", mock.Anything, []string{}).Return(nil)

	err := runConfigureTool(&linter, "--tool", "tool1")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestConfigure_Tool_Failed(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectConfigurator
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectConfigurator{},
	}
	linter.On("Configure", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runConfigureTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestConfigure_Tool_NotFound(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectConfigurator
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectConfigurator{},
	}

	err := runConfigureTool(&linter, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestConfigure_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
	}

	err := runConfigureTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support configuration", err.Error())
	}
}
