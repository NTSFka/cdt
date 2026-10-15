package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestPylint_DetectPylint(t *testing.T) {
	pylint := test.RunDetectToolLastFound(t, pylintDetect, []string{"pylint"})

	assert.Equal(t, tool.IdPylint, pylint.Id())
}

func TestPylint_DetectPylint_NotFound(t *testing.T) {
	pylint := test.RunDetectToolNotFound(t, pylintDetect, []string{"pylint"})

	assert.Equal(t, tool.IdPylint, pylint.Id())
}

func TestPylint_DetectPylint_Config(t *testing.T) {
	pylint := test.RunDetectToolConfig(t, pylintDetect, "pylint")

	assert.Equal(t, tool.IdPylint, pylint.Id())
}

func TestPylint_Pylint_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		pylintLinter,
		internal.ProjectLinterOptions{},
		[]string{"*"},
	)
}

func TestPylint_Pylint_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		pylintLinter,
		internal.ProjectLinterOptions{},
		[]string{"*"},
	)
}

func TestPylint_Pylint_LintFiles(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesSuccess(
		t,
		pylintLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestPylint_Pylint_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesFailed(
		t,
		pylintLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestPylint_Pylint_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		pylintLinter,
		internal.LintReportFormatRaw,
		[]string{"*"},
		nil,
	)
}

func TestPylint_Pylint_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	formats := []internal.LintReportFormat{
		internal.LintReportFormatJson,
		internal.LintReportFormatJUnit,
		internal.LintReportFormatGitHub,
		internal.LintReportFormatGitLab,
		internal.LintReportFormatTeamCity,
		"test-unsupported",
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeLintFilesOutputFormatUnsupported(
				t,
				pylintLinter,
				format,
				nil,
			)
		})
	}
}

func TestPylint_Pylint_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(
		t,
		pylintLinter,
		[]string{"*"},
		nil,
	)
}

func pylintDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPylint(ctx, options)
}

func pylintLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewPylint(executable)
}
