package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestNilAway_DetectNilAway(t *testing.T) {
	nilAway := test.RunDetectToolLastFound(t, nilAwayDetect, []string{"nilaway"})

	assert.Equal(t, tool.IdNilAway, nilAway.Id())
}

func TestNilAway_DetectNilAway_NotFound(t *testing.T) {
	nilAway := test.RunDetectToolNotFound(t, nilAwayDetect, []string{"nilaway"})

	assert.Equal(t, tool.IdNilAway, nilAway.Id())
}

func TestNilAway_DetectNilAway_Config(t *testing.T) {
	nilAway := test.RunDetectToolConfig(t, nilAwayDetect, "nilaway")

	assert.Equal(t, tool.IdNilAway, nilAway.Id())
}

func TestNilAway_NilAway_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		nilAwayLinter,
		internal.ProjectLinterOptions{},
		[]string{"./..."},
	)
}

func TestNilAway_NilAway_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		nilAwayLinter,
		internal.ProjectLinterOptions{},
		[]string{"./..."},
	)
}

func TestNilAway_NilAway_LintFiles(t *testing.T) {
	fileNames := []string{"mod1"}

	test.InvokeLintFilesSuccess(
		t,
		nilAwayLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestNilAway_NilAway_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"mod1"}

	test.InvokeLintFilesFailed(
		t,
		nilAwayLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestNilAway_NilAway_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		nilAwayLinter,
		".",
		internal.LintReportFormatRaw,
		[]string{"./..."},
		nil,
	)
}

func TestNilAway_NilAway_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.InvokeLintFilesOutputFormatUnsupported(t, nilAwayLinter, ".", format, nil)
		})
	}
}

func TestNilAway_NilAway_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(
		t,
		nilAwayLinter,
		".",
		[]string{"./..."},
		nil,
	)
}

func nilAwayDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectNilAway(ctx, options)
}

func nilAwayLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewNilAway(executable)
}
