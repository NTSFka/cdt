package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func runLint(linter internal.ProjectLinter, args ...string) error {
	return test.RunCommand(NewLintCommand(), internal.Context{
		Workflow: internal.Workflow{
			Linter: linter,
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
	linter := test.NewProjectLinter(t)
	linter.On("LintAll", mock.Anything, []string{}).
		Return(nil)

	err := runLint(linter)

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_LintAll_Failure(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On("LintAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runLint(linter)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

type testLinterTool struct {
	internal.ExecutableTool
	test.ProjectLinter
}

func newTestLinterTool(t *testing.T) *testLinterTool {
	linter := &testLinterTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectLinter{},
	}
	linter.Test(t)
	return linter
}

func TestLint_Tool_Success(t *testing.T) {
	linter := newTestLinterTool(t)
	linter.On("LintAll", mock.Anything, []string{}).
		Return(nil)

	err := runLintTool(linter, "--tool", "tool1")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_Tool_Failed(t *testing.T) {
	linter := newTestLinterTool(t)
	linter.On("LintAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runLintTool(linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestLint_Tool_NotFound(t *testing.T) {
	linter := newTestLinterTool(t)

	err := runLintTool(linter, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestLint_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runLintTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support linting", err.Error())
	}
}

func TestLint_LintFiles_Success(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).
		Return(nil)

	err := runLint(linter, "file1", "file2")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_LintFiles_Failure(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On("LintFiles", mock.Anything, []string{"file1", "file2"}, []string{}).
		Return(errors.New("failed"))

	err := runLint(linter, "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}
