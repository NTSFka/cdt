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

func TestNilAway_AnalyzeFiles_All(t *testing.T) {
	test.InvokeAnalyzeFilesSuccess(
		t,
		nilAwayAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"./..."},
	)
}

func TestNilAway_AnalyzeFiles_All_Failed(t *testing.T) {
	test.InvokeAnalyzeFilesFailed(
		t,
		nilAwayAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"./..."},
	)
}

func TestNilAway_AnalyzeFiles(t *testing.T) {
	fileNames := []string{"mod1"}

	test.InvokeAnalyzeFilesSuccess(
		t,
		nilAwayAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestNilAway_AnalyzeFiles_Failed(t *testing.T) {
	fileNames := []string{"mod1"}

	test.InvokeAnalyzeFilesFailed(
		t,
		nilAwayAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestNilAway_AnalyzeFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeAnalyzeFilesOutputFormatCheck(
		t,
		nilAwayAnalyzer,
		".",
		internal.AnalyzeReportFormatRaw,
		[]string{"./..."},
		nil,
	)
}

func TestNilAway_AnalyzeFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []internal.AnalyzeReportFormat{
		internal.AnalyzeReportFormatJson,
		internal.AnalyzeReportFormatJUnit,
		internal.AnalyzeReportFormatGitHub,
		internal.AnalyzeReportFormatGitLab,
		internal.AnalyzeReportFormatTeamCity,
		"test-unsupported",
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeAnalyzeFilesOutputFormatUnsupported(t, nilAwayAnalyzer, ".", format, nil)
		})
	}
}

func TestNilAway_AnalyzeFiles_OutputFile(t *testing.T) {
	test.InvokeAnalyzeFilesOutputFile(
		t,
		nilAwayAnalyzer,
		".",
		[]string{"./..."},
		nil,
	)
}

func nilAwayDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectNilAway(ctx, options)
}

func nilAwayAnalyzer(executable func() (*internal.Executable, error)) internal.ProjectAnalyzer {
	return tool.NewNilAway(executable)
}
