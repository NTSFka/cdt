package cli

import (
	"cdt/internal"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runFormat(formatter internal.ProjectFormatter, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "format")
	runArgs = append(runArgs, args...)

	return runMainWithWorkflow(internal.Workflow{
		Formatter: formatter,
	}, runArgs...)
}

func TestFormatConfigDefault(t *testing.T) {
	config := runMainGetConfig("format")

	assert.Equal(t, ".", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestFormatConfigCustomRootDirectory(t *testing.T) {
	config := runMainGetConfig("format", "--directory", "data/test")

	assert.Equal(t, "data/test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestFormatConfigCustomRootDirectoryShort(t *testing.T) {
	config := runMainGetConfig("format", "-d", "data/short-test")

	assert.Equal(t, "data/short-test", config.RootDirectory)
	assert.Nil(t, config.BuildDirectory)
}

func TestFormatCannotBeFormatted(t *testing.T) {
	err := runFormat(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support source formatting", err.Error())
	}
}

func TestFormatAllSuccess(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatAll", mock.Anything, []string{}).Return(nil)

	err := runFormat(&formatter)

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormatAllFailure(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormatFilesSuccess(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(nil)

	err := runFormat(&formatter, "file1", "file2")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormatFilesFailure(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter, "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormatCheckAllSuccess(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatCheckAll", mock.Anything, []string{}).Return(nil)

	err := runFormat(&formatter, "--check")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormatCheckAllFailure(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatCheckAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter, "--check")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}

func TestFormatCheckFilesSuccess(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatCheckFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(nil)

	err := runFormat(&formatter, "--check", "file1", "file2")

	assert.NoError(t, err)
	formatter.AssertExpectations(t)
}

func TestFormatCheckFilesFailure(t *testing.T) {
	formatter := testProjectFormatter{}
	formatter.On("FormatCheckFiles", mock.Anything, []string{"file1", "file2"}, []string{}).Return(errors.New("failed"))

	err := runFormat(&formatter, "--check", "file1", "file2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	formatter.AssertExpectations(t)
}
