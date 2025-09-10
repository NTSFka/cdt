package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func runConfigure(ctx context.Context, configurator internal.ProjectConfigurator, args ...string) error {
	return test.RunCommand(ctx, NewConfigureCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Configurator: configurator,
			},
		},
	}, args...)
}

func runConfigureTool(ctx context.Context, configurator internal.Tool, args ...string) error {
	return test.RunCommand(ctx, NewConfigureCommand(), internal.Context{
		Tools: []internal.Tool{
			configurator,
		},
	}, args...)
}

func TestConfigure_NotSupported(t *testing.T) {
	err := runConfigure(context.Background(), nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support configuration", err.Error())
	}
}

func TestConfigure_Configure_Success(t *testing.T) {
	configurator := test.NewProjectConfigurator(t)
	configurator.On("Configure", mock.Anything, mock.Anything, []string{}).
		Return(nil)

	err := runConfigure(context.Background(), configurator)

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestConfigure_Configure_Failure(t *testing.T) {
	configurator := test.NewProjectConfigurator(t)
	configurator.On("Configure", mock.Anything, mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runConfigure(context.Background(), configurator)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
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
	configurator.On("Configure", mock.Anything, mock.Anything, []string{}).
		Return(nil)

	err := runConfigureTool(context.Background(), configurator, "--tool", "tool1")

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestConfigure_Tool_Failed(t *testing.T) {
	configurator := newTestConfiguratorTool(t)
	configurator.On("Configure", mock.Anything, mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runConfigureTool(context.Background(), configurator, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestConfigure_Tool_NotFound(t *testing.T) {
	configurator := newTestConfiguratorTool(t)

	err := runConfigureTool(context.Background(), configurator, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestConfigure_Tool_NotSupported(t *testing.T) {
	configurator := &struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runConfigureTool(context.Background(), configurator, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support configuration", err.Error())
	}
}
