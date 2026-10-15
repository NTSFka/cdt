package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestBlack_DetectBlack(t *testing.T) {
	black := test.RunDetectToolLastFound(t, blackDetect, []string{"black"})

	assert.Equal(t, tool.IdBlack, black.Id())
}

func TestBlack_DetectBlack_NotFound(t *testing.T) {
	black := test.RunDetectToolNotFound(t, blackDetect, []string{"black"})

	assert.Equal(t, tool.IdBlack, black.Id())
}

func TestBlack_DetectBandit_Config(t *testing.T) {
	black := test.RunDetectToolConfig(t, blackDetect, "black")

	assert.Equal(t, tool.IdBlack, black.Id())
}

func TestBlack_Black_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		blackFormatter,
		internal.ProjectFormatterOptions{},
		[]string{"."},
	)
}

func TestBlack_Black_FormatFiles_All_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		blackFormatter,
		internal.ProjectFormatterOptions{},
		[]string{"."},
	)
}

func TestBlack_Black_FormatFiles(t *testing.T) {
	fileNames := []string{"tests/*", "/path/to/file.py"}

	test.InvokeFormatFilesSuccess(
		t,
		blackFormatter,
		internal.ProjectFormatterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestBlack_Black_FormatFiles_Failed(t *testing.T) {
	fileNames := []string{"tests/*", "/path/to/file.py"}

	test.InvokeFormatFilesFailed(
		t,
		blackFormatter,
		internal.ProjectFormatterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestBlack_Black_FormatFiles_CheckAll(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		blackFormatter,
		internal.ProjectFormatterOptions{
			CheckOnly: true,
		},
		[]string{"--check", "."},
	)
}

func TestBlack_Black_FormatFiles_Check(t *testing.T) {
	fileNames := []string{"tests/*", "/path/to/file.py"}

	test.InvokeFormatFilesSuccess(
		t,
		blackFormatter,
		internal.ProjectFormatterOptions{
			Filenames: &fileNames,
			CheckOnly: true,
		},
		append([]string{"--check"}, fileNames...),
	)
}

func blackDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectBlack(ctx, options)
}

func blackFormatter(executable func() (*internal.Executable, error)) internal.ProjectFormatter {
	return tool.NewBlack(executable)
}
