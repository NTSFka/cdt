package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runLint(linter internal.ProjectLinter, args ...string) error {
	return test.RunCommand(NewLintCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Linter: linter,
			},
		},
	}, args...)
}

func TestLintCannotBeLinted(t *testing.T) {
	err := runLint(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support linting", err.Error())
	}
}

func TestLintAllSuccess(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintAll", mock.Anything, []string{}).Return(nil)

	err := runLint(&linter)

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLintAllFailure(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runLint(&linter)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestLintFilesSuccess(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(nil)

	err := runLint(&linter, "file1", "file2")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLintFilesFailure(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(errors.New("failed"))

	err := runLint(&linter, "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}
