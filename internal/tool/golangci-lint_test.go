package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	golangCILint := test.RunDetectToolLastFound(t, golangCILintDetect, []string{"golangci-lint"})

	assert.Equal(t, tool.IdGolangCILint, golangCILint.Id())
}

func TestGolangCILint_DetectGolangCILint_NotFound(t *testing.T) {
	golangCILint := test.RunDetectToolNotFound(t, golangCILintDetect, []string{"golangci-lint"})

	assert.Equal(t, tool.IdGolangCILint, golangCILint.Id())
}

func TestGolangCILint_DetectGolangCILint_Config(t *testing.T) {
	golangCILint := test.RunDetectToolConfig(t, golangCILintDetect, "golangci-lint")

	assert.Equal(t, tool.IdGolangCILint, golangCILint.Id())
}

func TestGolangCILint_GolangCILint_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		golangCILintLinter,
		internal.ProjectLinterOptions{},
		[]string{"run"},
	)
}

func TestGolangCILint_GolangCILint_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		golangCILintLinter,
		internal.ProjectLinterOptions{},
		[]string{"run"},
	)
}

func TestGolangCILint_GolangCILint_LintFiles(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		golangCILintLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"mod1"}},
		[]string{"run", "mod1"},
	)
}

func TestGolangCILint_GolangCILint_LintFiles_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		golangCILintLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"mod1"}},
		[]string{"run", "mod1"},
	)
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		golangCILintLinter,
		internal.LintReportFormatRaw,
		[]string{"run"},
		nil,
	)
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.InvokeLintFilesOutputFormatUnsupported(t, golangCILintLinter, format, nil)
		})
	}
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(
		t,
		golangCILintLinter,
		[]string{"run"},
		nil,
	)
}

func TestGolangCILint_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"fmt"},
	)
}

func TestGolangCILint_FormatFiles_All_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"fmt"},
	)
}

func TestGolangCILint_FormatFiles(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Filenames:   &[]string{"file1"},
		},
		[]string{"fmt", "file1"},
	)
}

func TestGolangCILint_FormatFiles_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Filenames:   &[]string{"file1"},
		},
		[]string{"fmt", "file1"},
	)
}

func TestGolangCILint_FormatFiles_CheckAll(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
		},
		[]string{"fmt", "--diff"},
	)
}

func TestGolangCILint_FormatFiles_CheckAll_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
		},
		[]string{"fmt", "--diff"},
	)
}

func TestGolangCILint_FormatFiles_Check(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
			Filenames:   &[]string{"file1"},
		},
		[]string{"fmt", "--diff", "file1"},
	)
}

func TestGolangCILint_FormatFiles_Check_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		golangCILintFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
			Filenames:   &[]string{"file1"},
		},
		[]string{"fmt", "--diff", "file1"},
	)
}

func golangCILintDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectGolangCILint(ctx, options)
}

func golangCILintLinter(
	executable func() (*internal.Executable, error),
) internal.ProjectLinter {
	return tool.NewGolangCILint(executable)
}

func golangCILintFormatter(
	executable func() (*internal.Executable, error),
) internal.ProjectFormatter {
	return tool.NewGolangCILint(executable)
}
