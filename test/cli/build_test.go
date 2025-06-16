package cli

import (
	"cdt/internal"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runBuild(builder internal.ProjectBuilder, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "build")
	runArgs = append(runArgs, args...)

	return runMainWithWorkflow(internal.Workflow{
		Builder: builder,
	}, runArgs...)
}

func TestBuildConfigDefault(t *testing.T) {
	config := runMainGetConfig("build")

	assert.Equal(t, ".", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestBuildConfigCustomRootDirectory(t *testing.T) {
	config := runMainGetConfig("build", "--directory", "data/test")

	assert.Equal(t, "data/test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestBuildConfigCustomRootDirectoryShort(t *testing.T) {
	config := runMainGetConfig("build", "-d", "data/short-test")

	assert.Equal(t, "data/short-test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestBuildConfigCustomBuildDirectory(t *testing.T) {
	config := runMainGetConfig("build", "--build", "data/test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/test", *config.BuildDirectory)
	}
}

func TestBuildConfigCustomBuildDirectoryShort(t *testing.T) {
	config := runMainGetConfig("build", "-b", "data/short-test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/short-test", *config.BuildDirectory)
	}
}

func TestBuildNotSupported(t *testing.T) {
	err := runBuild(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support building", err.Error())
	}
}

func TestBuildAllSuccess(t *testing.T) {
	builder := testProjectBuilder{}
	builder.On("BuildAll", mock.Anything, []string{}).Return(nil)

	err := runBuild(&builder)

	assert.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuildAllFailure(t *testing.T) {
	builder := testProjectBuilder{}
	builder.On("BuildAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runBuild(&builder)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	builder.AssertExpectations(t)
}

func TestBuildTargetsSuccess(t *testing.T) {
	builder := testProjectBuilder{}
	builder.On("BuildTargets", mock.Anything, []string{"target1", "target2"}, []string{}).Return(nil)

	err := runBuild(&builder, "target1", "target2")

	assert.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuildTargetsFailure(t *testing.T) {
	builder := testProjectBuilder{}
	builder.On("BuildTargets", mock.Anything, []string{"target1", "target2"}, []string{}).Return(errors.New("failed"))

	err := runBuild(&builder, "target1", "target2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	builder.AssertExpectations(t)
}
