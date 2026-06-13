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
	test.RunLintFilesSuccess(
		t,
		nilAwayBuildLinter,
		internal.ProjectLinterOptions{},
		[]string{"./..."},
	)
}

func TestNilAway_NilAway_LintFiles(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		nilAwayBuildLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"mod1"}},
		[]string{"mod1"},
	)
}

func TestNilAway_NilAway_LintFiles_OutputFormat(t *testing.T) {
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
				nilAwayBuildLinter,
				data.format,
				append([]string{"./..."}, data.args...),
				nil,
			)
		})
	}
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
			test.RunLintFilesOutputFormatUnsupported(t, nilAwayBuildLinter, format, nil)
		})
	}
}

func TestNilAway_NilAway_LintFiles_OutputFile(t *testing.T) {
	test.RunLintFilesOutputFile(
		t,
		nilAwayBuildLinter,
		[]string{"./..."},
		nil,
	)
}

func nilAwayDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectNilAway(ctx, options)
}

func nilAwayBuildLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewNilAway(executable)
}
