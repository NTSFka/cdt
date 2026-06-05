package command_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/command"
	"cdt/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func runLint(ctx context.Context, linter internal.ProjectLinter, args ...string) error {
	return test.RunCommand(ctx, command.NewLintCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Linter: linter,
			},
		},
	}, args...)
}

func runLintTool(ctx context.Context, linter internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewLintCommand(), internal.Context{
		Tools: []internal.Tool{
			linter,
		},
	}, args...)
}

func TestLint_Lint_CannotBeLinted(t *testing.T) {
	err := runLint(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support linting", err.Error())
}

func TestLint_LintFiles_All_Success(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On("LintFiles", mock.Anything, mock.MatchedBy(func(options internal.ProjectLinterOptions) bool {
		return assert.Nil(t, options.Filenames)
	})).
		Return(nil)

	err := runLint(t.Context(), linter)

	require.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_LintFiles_All_Failure(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On("LintFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runLint(t.Context(), linter)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

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
	linter.On("LintFiles", mock.Anything, mock.Anything).
		Return(nil)

	err := runLintTool(t.Context(), linter, "--tool", "tool1")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_Tool_Failed(t *testing.T) {
	linter := newTestLinterTool(t)
	linter.On("LintFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runLintTool(t.Context(), linter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	linter.AssertExpectations(t)
}

func TestLint_Tool_NotFound(t *testing.T) {
	linter := newTestLinterTool(t)

	err := runLintTool(t.Context(), linter, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())

	linter.AssertExpectations(t)
}

func TestLint_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runLintTool(t.Context(), &linter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support linting", err.Error())
}

func TestLint_LintFiles_Success(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On(
		"LintFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectLinterOptions) bool {
			return assert.NotNil(t, opts.Filenames) &&
				assert.ElementsMatch(t, []string{"file1", "file2"}, *opts.Filenames)
		}),
	).
		Return(nil)

	err := runLint(t.Context(), linter, "file1", "file2")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_LintFiles_Failure(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On(
		"LintFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectLinterOptions) bool {
			return assert.NotNil(t, opts.Filenames) &&
				assert.ElementsMatch(t, []string{"file1", "file2"}, *opts.Filenames)
		}),
	).
		Return(errors.New("failed"))

	err := runLint(t.Context(), linter, "file1", "file2")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	linter.AssertExpectations(t)
}
func TestLint_LintFiles_CustomOutput_Raw(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On("LintFiles", mock.Anything, mock.MatchedBy(func(options internal.ProjectLinterOptions) bool {
		return assert.Equal(t, internal.LintReportFormatRaw, options.Output.Format) &&
			assert.Nil(t, options.Output.Filename)
	})).
		Return(nil)

	err := runLint(t.Context(), linter, "--output", "raw")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestLint_LintFiles_CustomOutput_Json(t *testing.T) {
	linter := test.NewProjectLinter(t)
	linter.On("LintFiles", mock.Anything, mock.MatchedBy(func(options internal.ProjectLinterOptions) bool {
		return assert.Equal(t, internal.LintReportFormatJson, options.Output.Format) &&
			assert.NotNil(t, options.Output.Filename) &&
			assert.Equal(t, "lint.json", *options.Output.Filename)
	})).
		Return(nil)

	err := runLint(t.Context(), linter, "-o", "json:lint.json")

	require.NoError(t, err)
	linter.AssertExpectations(t)
}
