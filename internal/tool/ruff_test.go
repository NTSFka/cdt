package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestRuff_DetectRuff(t *testing.T) {
	ruff := test.RunDetectToolLastFound(t, ruffDetect, []string{"ruff"})

	assert.Equal(t, tool.IdRuff, ruff.Id())
}

func TestRuff_DetectRuff_NotFound(t *testing.T) {
	ruff := test.RunDetectToolNotFound(t, ruffDetect, []string{"ruff"})

	assert.Equal(t, tool.IdRuff, ruff.Id())
}

func TestRuff_DetectRuff_Config(t *testing.T) {
	ruff := test.RunDetectToolConfig(t, ruffDetect, "ruff")

	assert.Equal(t, tool.IdRuff, ruff.Id())
}

func TestRuff_Ruff_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		ruffLinter,
		internal.ProjectLinterOptions{},
		[]string{"check"},
	)
}

func TestRuff_Ruff_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		ruffLinter,
		internal.ProjectLinterOptions{},
		[]string{"check"},
	)
}

func TestRuff_Ruff_LintFiles(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesSuccess(
		t,
		ruffLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"check"}, fileNames...),
	)
}

func TestRuff_Ruff_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesFailed(
		t,
		ruffLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"check"}, fileNames...),
	)
}

func TestRuff_Ruff_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		ruffLinter,
		".",
		internal.LintReportFormatRaw,
		[]string{"check"},
		nil,
	)
}

func TestRuff_Ruff_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.InvokeLintFilesOutputFormatUnsupported(t, ruffLinter, ".", format, nil)
		})
	}
}

func TestRuff_Ruff_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(
		t,
		ruffLinter,
		".",
		[]string{"check"},
		nil,
	)
}

func TestRuff_Ruff_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		ruffFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"format"},
	)
}

func TestRuff_Ruff_FormatFiles(t *testing.T) {
	fileNames := []string{"tests/*"}

	test.InvokeFormatFilesSuccess(
		t,
		ruffFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Filenames:   &fileNames,
		},
		append([]string{"format"}, fileNames...),
	)
}

func TestRuff_Ruff_FormatFiles_CheckAll(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		ruffFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
		},
		[]string{"format", "--check"},
	)
}

func TestRuff_Ruff_FormatFiles_Check(t *testing.T) {
	fileNames := []string{"tests/*", "/path/to/file.py"}

	test.InvokeFormatFilesSuccess(
		t,
		ruffFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
			Filenames:   &fileNames,
		},
		append([]string{"format", "--check"}, fileNames...),
	)
}

func ruffDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectRuff(ctx, options)
}

func ruffLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewRuff(executable)
}

func ruffFormatter(executable func() (*internal.Executable, error)) internal.ProjectFormatter {
	return tool.NewRuff(executable)
}
