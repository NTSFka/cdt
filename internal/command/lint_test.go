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

func runLintTool(linter internal.Tool, args ...string) error {
	return test.RunCommand(NewLintCommand(), internal.Context{
		Project: internal.Project{},
		Tools: []internal.Tool{
			linter,
		},
	}, args...)
}

func TestLint_Lint_CannotBeLinted(t *testing.T) {
	err := runLint(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support linting", err.Error())
	}
}

func TestLint_LintAll_Success(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintAll", mock.Anything, []string{}).Return(nil)

	err := runLint(&linter)

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_LintAll_Failure(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runLint(&linter)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestLint_Tool_Success(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectLinter{},
	}
	linter.On("LintAll", mock.Anything, []string{}).Return(nil)

	err := runLintTool(&linter, "--tool", "tool1")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_Tool_Failed(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectLinter{},
	}
	linter.On("LintAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runLintTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestLint_Tool_NotFound(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectLinter{},
	}

	err := runLintTool(&linter, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestLint_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
	}

	err := runLintTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support linting", err.Error())
	}
}

func TestLint_LintFiles_Success(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(nil)

	err := runLint(&linter, "file1", "file2")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_LintFiles_Failure(t *testing.T) {
	linter := test.ProjectLinter{}
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(errors.New("failed"))

	err := runLint(&linter, "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}
