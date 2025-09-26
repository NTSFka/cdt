package command_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/command"
	"cdt/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func runConfigure(
	ctx context.Context,
	configurator internal.ProjectConfigurator,
	args ...string,
) error {
	return test.RunCommand(ctx, command.NewConfigureCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Configurator: configurator,
			},
		},
	}, args...)
}

func runConfigureTool(ctx context.Context, configurator internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewConfigureCommand(), internal.Context{
		Tools: []internal.Tool{
			configurator,
		},
	}, args...)
}

func TestConfigure_NotSupported(t *testing.T) {
	err := runConfigure(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support configuration", err.Error())
}

func TestConfigure_Configure_Success(t *testing.T) {
	configurator := test.NewProjectConfigurator(t)
	configurator.On("Configure", mock.Anything, mock.Anything).
		Return(nil)

	err := runConfigure(t.Context(), configurator)

	require.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestConfigure_Configure_Failure(t *testing.T) {
	configurator := test.NewProjectConfigurator(t)
	configurator.On("Configure", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runConfigure(t.Context(), configurator)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())
	configurator.AssertExpectations(t)
}

type newConfiguratorTool struct {
	internal.ExecutableTool
	test.ProjectConfigurator
}

func newTestConfiguratorTool(t *testing.T) *newConfiguratorTool {
	tool := newConfiguratorTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectConfigurator{},
	}
	tool.Test(t)

	return &tool
}

func TestConfigure_Tool_Success(t *testing.T) {
	configurator := newTestConfiguratorTool(t)
	configurator.On("Configure", mock.Anything, mock.Anything).
		Return(nil)

	err := runConfigureTool(t.Context(), configurator, "--tool", "tool1")

	require.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestConfigure_Tool_Failed(t *testing.T) {
	configurator := newTestConfiguratorTool(t)
	configurator.On("Configure", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runConfigureTool(t.Context(), configurator, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())
	configurator.AssertExpectations(t)
}

func TestConfigure_Tool_NotFound(t *testing.T) {
	configurator := newTestConfiguratorTool(t)

	err := runConfigureTool(t.Context(), configurator, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())
	configurator.AssertExpectations(t)
}

func TestConfigure_Tool_NotSupported(t *testing.T) {
	configurator := &struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runConfigureTool(t.Context(), configurator, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support configuration", err.Error())
}
