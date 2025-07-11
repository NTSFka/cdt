package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runWithEnvironment(environment internal.Environment, args ...string) error {
	return test.RunCommand(EnvironmentCommand, internal.Context{
		Environment: environment,
	}, args...)
}

func TestEnvironment_Status_Running(t *testing.T) {
	env := test.Environment{}

	env.On("Id").Return("test")
	env.On("IsRunning", mock.Anything).Return(true)

	err := runWithEnvironment(&env, "status")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Status_NotRunning(t *testing.T) {
	env := test.Environment{}

	env.On("Id").Return("test")
	env.On("IsRunning", mock.Anything).Return(false)

	err := runWithEnvironment(&env, "status")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Start(t *testing.T) {
	env := test.Environment{}

	env.On("Start", mock.Anything).Return(nil)

	err := runWithEnvironment(&env, "start")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Start_Failed(t *testing.T) {
	env := test.Environment{}

	env.On("Start", mock.Anything).Return(errors.New("failed"))

	err := runWithEnvironment(&env, "start")

	assert.EqualError(t, err, "environment start failed: failed")
	env.AssertExpectations(t)
}

func TestEnvironment_Stop(t *testing.T) {
	env := test.Environment{}

	env.On("Stop", mock.Anything).Return(nil)

	err := runWithEnvironment(&env, "stop")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Stop_Failed(t *testing.T) {
	env := test.Environment{}

	env.On("Stop", mock.Anything).Return(errors.New("failed"))

	err := runWithEnvironment(&env, "stop")

	assert.EqualError(t, err, "environment stop failed: failed")
	env.AssertExpectations(t)
}
