package cli

import (
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvironment_ConfigDefault(t *testing.T) {
	config := runMainGetConfig("environment")

	assert.Equal(t, ".", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
	assert.Nil(t, config.Environment)
}

func TestEnvironment_ConfigCustomRootDirectory(t *testing.T) {
	config := runMainGetConfig("environment", "--directory", "data/test")

	assert.Equal(t, "data/test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
	assert.Nil(t, config.Environment)
}

func TestEnvironment_ConfigCustomRootDirectoryShort(t *testing.T) {
	config := runMainGetConfig("environment", "-d", "data/short-test")

	assert.Equal(t, "data/short-test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
	assert.Nil(t, config.Environment)
}

func TestEnvironment_ConfigCustomEnvironment(t *testing.T) {
	config := runMainGetConfig("environment", "--environment", "env:image")

	if assert.NotNil(t, config.Environment) {
		assert.Equal(t, "env", config.Environment.ToolName)
		assert.Equal(t, "image", config.Environment.Argument)
	}
}

func TestEnvironment_ConfigCustomEnvironmentShort(t *testing.T) {
	config := runMainGetConfig("environment", "-e", "env:image")

	if assert.NotNil(t, config.Environment) {
		assert.Equal(t, "env", config.Environment.ToolName)
		assert.Equal(t, "image", config.Environment.Argument)
	}
}

func TestEnvironment_Status_Running(t *testing.T) {
	env := test.Environment{}

	env.On("Id").Return("test")
	env.On("IsRunning").Return(true)
	env.On("Cleanup").Return(nil)

	err := runMainWithEnvironment(&env, "environment", "status")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Status_NotRunning(t *testing.T) {
	env := test.Environment{}

	env.On("Id").Return("test")
	env.On("IsRunning").Return(false)
	env.On("Cleanup").Return(nil)

	err := runMainWithEnvironment(&env, "environment", "status")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Start(t *testing.T) {
	env := test.Environment{}

	env.On("Start").Return(nil)
	env.On("Cleanup").Return(nil)

	err := runMainWithEnvironment(&env, "environment", "start")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Start_Failed(t *testing.T) {
	env := test.Environment{}

	env.On("Start").Return(errors.New("failed"))
	env.On("Cleanup").Return(nil)

	err := runMainWithEnvironment(&env, "environment", "start")

	assert.EqualError(t, err, "environment start failed: failed")
	env.AssertExpectations(t)
}

func TestEnvironment_Stop(t *testing.T) {
	env := test.Environment{}

	env.On("Stop").Return(nil)
	env.On("Cleanup").Return(nil)

	err := runMainWithEnvironment(&env, "environment", "stop")

	assert.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Stop_Failed(t *testing.T) {
	env := test.Environment{}

	env.On("Stop").Return(errors.New("failed"))
	env.On("Cleanup").Return(nil)

	err := runMainWithEnvironment(&env, "environment", "stop")

	assert.EqualError(t, err, "environment stop failed: failed")
	env.AssertExpectations(t)
}
