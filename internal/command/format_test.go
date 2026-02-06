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

func runFormat(ctx context.Context, formatter internal.ProjectFormatter, args ...string) error {
	return test.RunCommand(ctx, command.NewFormatCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Formatter: formatter,
			},
		},
	}, args...)
}

func runFormatTool(ctx context.Context, formatter internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewFormatCommand(), internal.Context{
		Tools: []internal.Tool{
			formatter,
		},
	}, args...)
}

func TestFormat_CannotBeFormatted(t *testing.T) {
	err := runFormat(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support source formatting", err.Error())
}

func TestFormat_FormatFiles_All_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On(
		"FormatFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectFormatterOptions) bool {
			return assert.False(t, opts.CheckOnly) && assert.Nil(t, opts.Filenames)
		}),
	).
		Return(nil)

	err := runFormat(t.Context(), formatter)

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_All_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runFormat(t.Context(), formatter)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}

type testFormatterTool struct {
	internal.ExecutableTool
	test.ProjectFormatter
}

func newFormatterTool(t *testing.T) *testFormatterTool {
	formatter := &testFormatterTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectFormatter{},
	}
	formatter.Test(t)

	return formatter
}

func TestFormat_Tool_Success(t *testing.T) {
	formatter := newFormatterTool(t)
	formatter.On("FormatFiles", mock.Anything, mock.Anything).
		Return(nil)

	err := runFormatTool(t.Context(), formatter, "--tool", "tool1")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_Tool_Failed(t *testing.T) {
	formatter := newFormatterTool(t)
	formatter.On("FormatFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runFormatTool(t.Context(), formatter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}

func TestFormat_Tool_NotFound(t *testing.T) {
	formatter := newFormatterTool(t)

	err := runFormatTool(t.Context(), formatter, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())

	formatter.AssertExpectations(t)
}

func TestFormat_Tool_NotSupported(t *testing.T) {
	formatter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runFormatTool(t.Context(), &formatter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support formatting", err.Error())
}

func TestFormat_FormatFiles_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On(
		"FormatFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectFormatterOptions) bool {
			return assert.NotNil(t, opts.Filenames) &&
				assert.Equal(t, []string{"file1", "file2"}, *opts.Filenames)
		}),
	).
		Return(nil)

	err := runFormat(t.Context(), formatter, "file1", "file2")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On(
		"FormatFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectFormatterOptions) bool {
			return assert.NotNil(t, opts.Filenames) &&
				assert.Equal(t, []string{"file1", "file2"}, *opts.Filenames)
		}),
	).
		Return(errors.New("failed"))

	err := runFormat(t.Context(), formatter, "file1", "file2")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_CheckAll_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On(
		"FormatFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectFormatterOptions) bool {
			return assert.True(t, opts.CheckOnly)
		}),
	).
		Return(nil)

	err := runFormat(t.Context(), formatter, "--check")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_CheckAll_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runFormat(t.Context(), formatter, "--check")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_Check_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On(
		"FormatFiles",
		mock.Anything,
		mock.MatchedBy(func(opts internal.ProjectFormatterOptions) bool {
			return assert.NotNil(t, opts.Filenames) &&
				assert.Equal(t, []string{"file1", "file2"}, *opts.Filenames) &&
				assert.True(t, opts.CheckOnly)
		}),
	).
		Return(nil)

	err := runFormat(t.Context(), formatter, "--check", "file1", "file2")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_Check_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatFiles", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runFormat(t.Context(), formatter, "--check", "file1", "file2")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}
