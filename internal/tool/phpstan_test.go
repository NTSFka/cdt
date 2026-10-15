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

func TestPHPStan_PHPStan_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		phpStanLinter,
		internal.ProjectLinterOptions{},
		[]string{"analyse"},
	)
}

func TestPHPStan_PHPStan_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		phpStanLinter,
		internal.ProjectLinterOptions{},
		[]string{"analyse"},
	)
}

func TestPHPStan_PHPStan_LintFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesSuccess(
		t,
		phpStanLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"analyse"}, fileNames...),
	)
}

func TestPHPStan_PHPStan_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesFailed(
		t,
		phpStanLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"analyse"}, fileNames...),
	)
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat(t *testing.T) {
	dataSet := []struct {
		format internal.LintReportFormat
		args   []string
	}{
		{internal.LintReportFormatRaw, []string{}},
		{internal.LintReportFormatJson, []string{"--error-format=json"}},
		{internal.LintReportFormatJUnit, []string{"--error-format=junit"}},
		{internal.LintReportFormatGitHub, []string{"--error-format=github"}},
		{internal.LintReportFormatGitLab, []string{"--error-format=gitlab"}},
		{internal.LintReportFormatTeamCity, []string{"--error-format=teamcity"}},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			test.InvokeLintFilesOutputFormatCheck(
				t,
				phpStanLinter,
				data.format,
				append([]string{"analyse"}, data.args...),
				nil,
			)
		})
	}
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	formats := []internal.LintReportFormat{
		"test-unsupported",
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeLintFilesOutputFormatUnsupported(
				t,
				phpStanLinter,
				format,
				nil,
			)
		})
	}
}

func TestPHPStan_PHPStan_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(
		t,
		phpStanLinter,
		[]string{"analyse"},
		nil,
	)
}

func phpStanDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHPStan(ctx, options)
}

func phpStanLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewPHPStan(executable)
}
