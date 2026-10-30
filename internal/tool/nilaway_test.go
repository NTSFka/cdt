package tool_test

import (
	"bytes"
	"cdt/internal/reporter"
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNilAway_AnalyzeFiles_Reporter_NoErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		nilAwayAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"-json", "./..."},
		"{}",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestNilAway_AnalyzeFiles_Reporter_SomeErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		nilAwayAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"-json", "./..."},
		`{
			"cdt/internal/command": {
				"nilaway": [
					{
						"posn": "test.go:25:13",
						"end": "test.go:25:13",
						"message": "some error"
					}
				]
			}
		}`,
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.NotEmpty(t, result.String())
}

func TestNilAway_AnalyzeFiles_Reporter_NoOutput(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		nilAwayAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"-json", "./..."},
		"",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestNilAway_AnalyzeFiles_Reporter_InvalidJson(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		nilAwayAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"-json", "./..."},
		".",
		func(err error) {
			require.ErrorContains(t, err, "invalid character '.'")
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
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
