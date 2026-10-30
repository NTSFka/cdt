package tool_test

import (
	"bytes"
	"cdt/internal/reporter"
	"cdt/internal/test"
	"context"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestPHPStan_AnalyzeFiles_All(t *testing.T) {
	test.InvokeAnalyzeFilesSuccess(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"analyse"},
	)
}

func TestPHPStan_AnalyzeFiles_All_Failed(t *testing.T) {
	test.InvokeAnalyzeFilesFailed(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"analyse"},
	)
}

func TestPHPStan_AnalyzeFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesSuccess(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		append([]string{"analyse"}, fileNames...),
	)
}

func TestPHPStan_AnalyzeFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesFailed(
		t,
		phpStanAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		append([]string{"analyse"}, fileNames...),
	)
}

func TestPHPStan_AnalyzeFiles_Reporter_NoErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		phpStanAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyse", "--error-format=json", "--quiet"}, "{}",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestPHPStan_AnalyzeFiles_Reporter_SomeErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		phpStanAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyse", "--error-format=json", "--quiet"},
		`{
			"totals": {},
			"files": {
				"file1.php": {
					"messages": [
						{"message": "test", "line": 1, "identifier": "undefined"}
					]
				}
			}
  		}`,
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.NotEmpty(t, result.String())
}

func TestPHPStan_AnalyzeFiles_Reporter_NoOutput(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		phpStanAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyse", "--error-format=json", "--quiet"},
		"",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestPHPStan_AnalyzeFiles_Reporter_InvalidJson(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		phpStanAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyse", "--error-format=json", "--quiet"},
		".",
		func(err error) {
			require.ErrorContains(t, err, "invalid character '.'")
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestPHPStan_AnalyzeFiles_OutputFile(t *testing.T) {
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
