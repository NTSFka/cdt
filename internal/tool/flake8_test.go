package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestFlake8_DetectFlake8(t *testing.T) {
	flake8 := test.RunDetectToolLastFound(t, flake8Detect, []string{"flake8"})

	assert.Equal(t, tool.IdFlake8, flake8.Id())
}

func TestFlake8_DetectFlake8_NotFound(t *testing.T) {
	flake8 := test.RunDetectToolNotFound(t, flake8Detect, []string{"flake8"})

	assert.Equal(t, tool.IdFlake8, flake8.Id())
}

func TestFlake8_DetectFlake8_Config(t *testing.T) {
	flake8 := test.RunDetectToolConfig(t, flake8Detect, "flake8")

	assert.Equal(t, tool.IdFlake8, flake8.Id())
}

func TestFlake8_Flake8_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		flake8Linter,
		internal.ProjectLinterOptions{},
		[]string{},
	)
}

func TestFlake8_Flake8_LintFiles(t *testing.T) {
	fileNames := []string{"file1.py", "/path/to/file2.py"}

	test.InvokeLintFilesSuccess(
		t,
		flake8Linter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestFlake8_Flake8_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesFailed(
		t,
		flake8Linter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestFlake8_Flake8_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		flake8Linter,
		internal.LintReportFormatRaw,
		[]string{},
		nil,
	)
}

func TestFlake8_Flake8_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []internal.LintReportFormat{
		internal.LintReportFormatJson,
		internal.LintReportFormatJUnit,
		internal.LintReportFormatGitHub,
		internal.LintReportFormatGitLab,
		internal.LintReportFormatTeamCity,
		"test-unsupported",
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeLintFilesOutputFormatUnsupported(t, flake8Linter, format, nil)
		})
	}
}

func TestFlake8_Flake8_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(t, flake8Linter, []string{}, nil)
}

func flake8Detect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectFlake8(ctx, options)
}

func flake8Linter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewFlake8(executable)
}
