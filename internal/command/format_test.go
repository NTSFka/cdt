package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func runFormat(formatter internal.ProjectFormatter, args ...string) error {
	return test.RunCommand(NewFormatCommand(), internal.Context{
		Workflow: internal.Workflow{
			Formatter: formatter,
		},
	}, args...)
}

func runFormatTool(formatter internal.Tool, args ...string) error {
	return test.RunCommand(NewFormatCommand(), internal.Context{
		Project: internal.Project{},
		Tools: []internal.Tool{
			formatter,
		},
	}, args...)
}

func TestFormat_CannotBeFormatted(t *testing.T) {
	err := runFormat(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support source formatting", err.Error())
	}
}

func TestFormat_FormatAll_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatAll", mock.Anything, []string{}).
		Return(nil)

	err := runFormat(formatter)

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatAll_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runFormat(formatter)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
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
	formatter.On("FormatAll", mock.Anything, []string{}).
		Return(nil)

	err := runFormatTool(formatter, "--tool", "tool1")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_Tool_Failed(t *testing.T) {
	formatter := newFormatterTool(t)
	formatter.On("FormatAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runFormatTool(formatter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormat_Tool_NotFound(t *testing.T) {
	formatter := newFormatterTool(t)

	err := runFormatTool(formatter, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormat_Tool_NotSupported(t *testing.T) {
	formatter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runFormatTool(&formatter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support formatting", err.Error())
	}
}

func TestFormat_FormatFiles_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatFiles", mock.Anything, []string{"file1", "file2"}, []string{}).
		Return(nil)

	err := runFormat(formatter, "file1", "file2")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatFiles", mock.Anything, []string{"file1", "file2"}, []string{}).
		Return(errors.New("failed"))

	err := runFormat(formatter, "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckAll_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckAll", mock.Anything, []string{}).
		Return(nil)

	err := runFormat(formatter, "--check")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckAll_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runFormat(formatter, "--check")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckFiles_Success(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckFiles", mock.Anything, []string{"file1", "file2"}, []string{}).
		Return(nil)

	err := runFormat(formatter, "--check", "file1", "file2")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckFiles_Failure(t *testing.T) {
	formatter := test.NewProjectFormatter(t)
	formatter.On("FormatCheckFiles", mock.Anything, []string{"file1", "file2"}, []string{}).
		Return(errors.New("failed"))

	err := runFormat(formatter, "--check", "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}
