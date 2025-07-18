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

func TestConfigureNotSupported(t *testing.T) {
	err := runConfigure(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support configuration", err.Error())
	}
}

func TestConfigureSuccess(t *testing.T) {
	configurator := test.ProjectConfigurator{}
	configurator.On("Configure", mock.Anything, []string{}).Return(nil)

	err := runConfigure(&configurator)

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestConfigureFailure(t *testing.T) {
	configurator := test.ProjectConfigurator{}
	configurator.On("Configure", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runConfigure(&configurator)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	configurator.AssertExpectations(t)
}
