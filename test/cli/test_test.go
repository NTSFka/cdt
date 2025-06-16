package cli

import (
	"cdt/internal"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runTest(tester internal.ProjectTester, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "test")
	runArgs = append(runArgs, args...)

	return runMainWithWorkflow(internal.Workflow{
		Tester: tester,
	}, runArgs...)
}

func TestTestConfigDefault(t *testing.T) {
	config := runMainGetConfig("test")

	assert.Equal(t, ".", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestTestConfigCustomRootDirectory(t *testing.T) {
	config := runMainGetConfig("test", "--directory", "data/test")

	assert.Equal(t, "data/test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestTestConfigCustomRootDirectoryShort(t *testing.T) {
	config := runMainGetConfig("test", "-d", "data/short-test")

	assert.Equal(t, "data/short-test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestTestConfigCustomBuildDirectory(t *testing.T) {
	config := runMainGetConfig("test", "--build", "data/test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/test", *config.BuildDirectory)
	}
}

func TestTestConfigCustomBuildDirectoryShort(t *testing.T) {
	config := runMainGetConfig("test", "-b", "data/short-test")

	assert.Equal(t, ".", config.RootDirectory)

	if assert.NotNil(t, config.BuildDirectory) {
		assert.Equal(t, "data/short-test", *config.BuildDirectory)
	}
}

func TestTestCannotBeTested(t *testing.T) {
	err := runTest(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support testing", err.Error())
	}
}

func TestTestAllSuccess(t *testing.T) {
	tester := testProjectTester{}
	tester.On("TestAll", mock.Anything, []string{}).Return(nil)

	err := runTest(&tester)

	assert.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTestAllFailure(t *testing.T) {
	tester := testProjectTester{}
	tester.On("TestAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runTest(&tester)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	tester.AssertExpectations(t)
}

func TestTestTargetsSuccess(t *testing.T) {
	tester := testProjectTester{}
	tester.On("Test", mock.Anything, "pattern", []string{}).Return(nil)

	err := runTest(&tester, "pattern")

	assert.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTestTargetsFailure(t *testing.T) {
	tester := testProjectTester{}
	tester.On("Test", mock.Anything, "pattern", []string{}).Return(errors.New("failed"))

	err := runTest(&tester, "pattern")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	tester.AssertExpectations(t)
}
