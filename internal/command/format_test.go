package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func runFormat(ctx context.Context, formatter internal.ProjectFormatter, args ...string) error {
	return test.RunCommand(ctx, NewFormatCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Formatter: formatter,
			},
		},
	}, args...)
}

func runFormatTool(ctx context.Context, formatter internal.Tool, args ...string) error {
	return test.RunCommand(ctx, NewFormatCommand(), internal.Context{
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

func TestFormat_FormatAll_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatAll", mock.Anything, mock.Anything).
		Return(nil)

	err := runFormat(t.Context(), formatter)

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatAll_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatAll", mock.Anything, mock.Anything).
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
	formatter.On("FormatAll", mock.Anything, mock.Anything).
		Return(nil)

	err := runFormatTool(t.Context(), formatter, "--tool", "tool1")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_Tool_Failed(t *testing.T) {
	formatter := newFormatterTool(t)
	formatter.On("FormatAll", mock.Anything, mock.Anything).
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
	formatter.On("FormatFiles", mock.Anything, mock.Anything, []string{"file1", "file2"}).
		Return(nil)

	err := runFormat(t.Context(), formatter, "file1", "file2")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatFiles", mock.Anything, mock.Anything, []string{"file1", "file2"}).
		Return(errors.New("failed"))

	err := runFormat(t.Context(), formatter, "file1", "file2")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckAll_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckAll", mock.Anything, mock.Anything).
		Return(nil)

	err := runFormat(t.Context(), formatter, "--check")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckAll_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckAll", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runFormat(t.Context(), formatter, "--check")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckFiles_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckFiles", mock.Anything, mock.Anything, []string{"file1", "file2"}).
		Return(nil)

	err := runFormat(t.Context(), formatter, "--check", "file1", "file2")

	require.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckFiles_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckFiles", mock.Anything, mock.Anything, []string{"file1", "file2"}).
		Return(errors.New("failed"))

	err := runFormat(t.Context(), formatter, "--check", "file1", "file2")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	formatter.AssertExpectations(t)
}
