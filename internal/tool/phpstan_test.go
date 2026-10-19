package tool_test

import (
	"cdt/internal/test"
	"context"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestPHPStan_DetectPHPStan_Composer(t *testing.T) {
	phpStan := test.RunDetectToolLastFound(
		t,
		phpStanDetect,
		[]string{
			filepath.Join("vendor", "bin", "phpstan"),
		},
	)

	assert.Equal(t, tool.IdPHPStan, phpStan.Id())
}

func TestPHPStan_DetectPHPStan_System(t *testing.T) {
	phpStan := test.RunDetectToolLastFound(
		t,
		phpStanDetect,
		[]string{
			filepath.Join("vendor", "bin", "phpstan"),
			"phpstan",
		},
	)

	assert.Equal(t, tool.IdPHPStan, phpStan.Id())
}

func TestPHPStan_DetectPHPStan_NotFound(t *testing.T) {
	phpStan := test.RunDetectToolNotFound(
		t,
		phpStanDetect,
		[]string{
			filepath.Join("vendor", "bin", "phpstan"),
			"phpstan",
		},
	)

	assert.Equal(t, tool.IdPHPStan, phpStan.Id())
}

func TestPHPStan_DetectPHPStan_Config(t *testing.T) {
	phpStan := test.RunDetectToolConfig(t, phpStanDetect, "phpstan")

	assert.Equal(t, tool.IdPHPStan, phpStan.Id())
}

func TestPHPStan_PHPStan_AnalyzeFiles_All(t *testing.T) {
	test.InvokeAnalyzeFilesSuccess(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"analyse"},
	)
}

func TestPHPStan_PHPStan_AnalyzeFiles_All_Failed(t *testing.T) {
	test.InvokeAnalyzeFilesFailed(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"analyse"},
	)
}

func TestPHPStan_PHPStan_AnalyzeFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesSuccess(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		append([]string{"analyse"}, fileNames...),
	)
}

func TestPHPStan_PHPStan_AnalyzeFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesFailed(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		append([]string{"analyse"}, fileNames...),
	)
}

func TestPHPStan_PHPStan_AnalyzeFiles_OutputFormat(t *testing.T) {
	dataSet := []struct {
		format internal.AnalyzeReportFormat
		args   []string
	}{
		{internal.AnalyzeReportFormatRaw, []string{}},
		{internal.AnalyzeReportFormatJson, []string{"--error-format=json"}},
		{internal.AnalyzeReportFormatJUnit, []string{"--error-format=junit"}},
		{internal.AnalyzeReportFormatGitHub, []string{"--error-format=github"}},
		{internal.AnalyzeReportFormatGitLab, []string{"--error-format=gitlab"}},
		{internal.AnalyzeReportFormatTeamCity, []string{"--error-format=teamcity"}},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			test.InvokeAnalyzeFilesOutputFormatCheck(
				t,
				phpStanAnalyzer,
				".",
				data.format,
				append([]string{"analyse"}, data.args...),
				nil,
			)
		})
	}
}

func TestPHPStan_PHPStan_AnalyzeFiles_OutputFormat_Unsupported(t *testing.T) {
	formats := []internal.AnalyzeReportFormat{
		"test-unsupported",
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeAnalyzeFilesOutputFormatUnsupported(
				t,
				phpStanAnalyzer,
				".",
				format,
				nil,
			)
		})
	}
}

func TestPHPStan_PHPStan_AnalyzeFiles_OutputFile(t *testing.T) {
	test.InvokeAnalyzeFilesOutputFile(
		t,
		phpStanAnalyzer,
		".",
		[]string{"analyse"},
		nil,
	)
}

func phpStanDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHPStan(ctx, options)
}

func phpStanAnalyzer(executable func() (*internal.Executable, error)) internal.ProjectAnalyzer {
	return tool.NewPHPStan(executable)
}
