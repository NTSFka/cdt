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
	test.RunLintFilesSuccess(t, pylintBuildLinter, internal.ProjectLinterOptions{}, []string{"*"})
}

func TestPylint_Pylint_LintFiles(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		pylintBuildLinter,
		internal.ProjectLinterOptions{
			Filenames: &[]string{"file.py", "/path/to/file2.py"},
		},
		[]string{"file.py", "/path/to/file2.py"},
	)
}

func TestPylint_Pylint_LintFiles_OutputFormat(t *testing.T) {
	dataSet := []struct {
		format internal.LintReportFormat
		args   []string
	}{
		{internal.LintReportFormatRaw, []string{}},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			test.RunLintFilesOutputFormatCheck(
				t,
				pylintBuildLinter,
				data.format,
				append([]string{"*"}, data.args...),
				nil,
			)
		})
	}
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
			test.RunLintFilesOutputFormatUnsupported(
				t,
				pylintBuildLinter,
				format,
				nil,
			)
		})
	}
}

func TestPylint_Pylint_LintFiles_OutputFile(t *testing.T) {
	test.RunLintFilesOutputFile(
		t,
		pylintBuildLinter,
		[]string{"*"},
		nil,
	)
}

func pylintDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPylint(ctx, options)
}

func pylintBuildLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewPylint(executable)
}
