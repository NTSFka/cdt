package cli

import (
	"cdt/internal"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runLint(linter internal.ProjectLinter, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "lint")
	runArgs = append(runArgs, args...)

	return runMainWithWorkflow(internal.Workflow{
		Linter: linter,
	}, runArgs...)
}

func TestLintConfigDefault(t *testing.T) {
	config := runMainGetConfig("lint")

	assert.Equal(t, ".", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestLintConfigCustomRootDirectory(t *testing.T) {
	config := runMainGetConfig("lint", "--directory", "data/test")

	assert.Equal(t, "data/test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestLintConfigCustomRootDirectoryShort(t *testing.T) {
	config := runMainGetConfig("lint", "-d", "data/short-test")

	assert.Equal(t, "data/short-test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestLintCannotBeLinted(t *testing.T) {
	err := runLint(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support linting", err.Error())
	}
}

func TestLintAllSuccess(t *testing.T) {
	linter := testProjectLinter{}
	linter.On("LintAll", mock.Anything, []string{}).Return(nil)

	err := runLint(&linter)

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLintAllFailure(t *testing.T) {
	linter := testProjectLinter{}
	linter.On("LintAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runLint(&linter)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestLintFilesSuccess(t *testing.T) {
	linter := testProjectLinter{}
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(nil)

	err := runLint(&linter, "file1", "file2")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLintFilesFailure(t *testing.T) {
	linter := testProjectLinter{}
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(errors.New("failed"))

	err := runLint(&linter, "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}
