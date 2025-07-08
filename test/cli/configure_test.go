package cli

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runConfigure(configurator internal.ProjectConfigurator, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "configure")
	runArgs = append(runArgs, args...)

	return runMainWithWorkflow(internal.Workflow{
		Configurator: configurator,
	}, runArgs...)
}

func TestConfigureConfigDefault(t *testing.T) {
	config := runMainGetConfig("configure")

	assert.Equal(t, ".", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestConfigureConfigCustomRootDirectory(t *testing.T) {
	config := runMainGetConfig("configure", "--directory", "data/test")

	assert.Equal(t, "data/test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestConfigureConfigCustomRootDirectoryShort(t *testing.T) {
	config := runMainGetConfig("configure", "-d", "data/short-test")

	assert.Equal(t, "data/short-test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestConfigureConfigCustomBuildDirectory(t *testing.T) {
	config := runMainGetConfig("configure", "--build", "data/test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/test", *config.BuildDirectory)
	}
}

func TestConfigureConfigCustomBuildDirectoryShort(t *testing.T) {
	config := runMainGetConfig("configure", "-b", "data/short-test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/short-test", *config.BuildDirectory)
	}
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
