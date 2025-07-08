package cli

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runRun(runner internal.ProjectRunner, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "run")
	runArgs = append(runArgs, args...)

	return runMainWithWorkflow(internal.Workflow{
		Runner: runner,
	}, runArgs...)
}

func TestRunConfigDefault(t *testing.T) {
	config := runMainGetConfig("run")

	assert.Equal(t, ".", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestRunConfigCustomRootDirectory(t *testing.T) {
	config := runMainGetConfig("run", "--directory", "data/test")

	assert.Equal(t, "data/test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestRunConfigCustomRootDirectoryShort(t *testing.T) {
	config := runMainGetConfig("run", "-d", "data/short-test")

	assert.Equal(t, "data/short-test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestRunConfigCustomBuildDirectory(t *testing.T) {
	config := runMainGetConfig("run", "--build", "data/test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/test", *config.BuildDirectory)
	}
}

func TestRunConfigCustomBuildDirectoryShort(t *testing.T) {
	config := runMainGetConfig("run", "-b", "data/short-test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/short-test", *config.BuildDirectory)
	}
}

func TestRunNotSupported(t *testing.T) {
	err := runRun(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support run of target", err.Error())
	}
}

func TestRunAllSuccess(t *testing.T) {
	runner := test.ProjectRunner{}
	runner.On("RunTarget", mock.Anything, "target1", []string{}).Return(nil)

	err := runRun(&runner, "target1")

	assert.NoError(t, err)
	runner.AssertExpectations(t)
}

func TestRunAllFailure(t *testing.T) {
	runner := test.ProjectRunner{}
	runner.On("RunTarget", mock.Anything, "target1", []string{}).Return(errors.New("failed"))

	err := runRun(&runner, "target1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	runner.AssertExpectations(t)
}
