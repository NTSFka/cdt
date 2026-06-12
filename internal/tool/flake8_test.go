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
	test.RunLintFilesSuccess(
		t,
		flake8BuildLinter,
		internal.ProjectLinterOptions{},
		[]string{},
	)
}

func TestFlake8_Flake8_LintFiles(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		flake8BuildLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"file.py", "/path/to/file2.py"}},
		[]string{"file.py", "/path/to/file2.py"},
	)
}

func TestFlake8_Flake8_LintFiles_OutputFormat_Raw(t *testing.T) {
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
				flake8BuildLinter,
				data.format,
				append([]string{}, data.args...),
				nil,
			)
		})
	}
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
			test.RunLintFilesOutputFormatUnsupported(t, flake8BuildLinter, format, nil)
		})
	}
}

func TestFlake8_Flake8_LintFiles_OutputFile(t *testing.T) {
	test.RunLintFilesOutputFile(t, flake8BuildLinter, []string{}, nil)
}

func flake8Detect(ctx context.Context, options tool.DetectOptions) internal.Tool {
	return tool.DetectFlake8(ctx, options)
}

func flake8BuildLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewFlake8(executable)
}
