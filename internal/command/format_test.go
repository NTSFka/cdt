package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runFormat(formatter internal.ProjectFormatter, args ...string) error {
	return test.RunCommand(NewFormatCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Formatter: formatter,
			},
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
	formatter := test.ProjectFormatter{}
	formatter.On("FormatAll", mock.Anything, []string{}).Return(nil)

	err := runFormat(&formatter)

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatAll_Failure(t *testing.T) {
	formatter := test.ProjectFormatter{}
	formatter.On("FormatAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormat_Tool_Success(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectFormatter
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectFormatter{},
	}
	linter.On("FormatAll", mock.Anything, []string{}).Return(nil)

	err := runFormatTool(&linter, "--tool", "tool1")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestFormat_Tool_Failed(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectFormatter
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectFormatter{},
	}
	linter.On("FormatAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runFormatTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestFormat_Tool_NotFound(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectFormatter
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectFormatter{},
	}

	err := runFormatTool(&linter, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestFormat_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
	}

	err := runFormatTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support formatting", err.Error())
	}
}

func TestFormat_FormatFiles_Success(t *testing.T) {
	formatter := test.ProjectFormatter{}
	formatter.On("FormatFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(nil)

	err := runFormat(&formatter, "file1", "file2")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatFiles_Failure(t *testing.T) {
	formatter := test.ProjectFormatter{}
	formatter.On("FormatFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter, "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckAll_Success(t *testing.T) {
	formatter := test.ProjectFormatter{}
	formatter.On("FormatCheckAll", mock.Anything, []string{}).Return(nil)

	err := runFormat(&formatter, "--check")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckAll_Failure(t *testing.T) {
	formatter := test.ProjectFormatter{}
	formatter.On("FormatCheckAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter, "--check")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckFiles_Success(t *testing.T) {
	formatter := test.ProjectFormatter{}
	formatter.On("FormatCheckFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(nil)

	err := runFormat(&formatter, "--check", "file1", "file2")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormat_FormatCheckFiles_Failure(t *testing.T) {
	formatter := test.ProjectFormatter{}
	formatter.On("FormatCheckFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter, "--check", "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}
