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

func TestMago_DetectMago_Composer(t *testing.T) {
	mago := test.RunDetectToolLastFound(
		t,
		magoDetect,
		[]string{
			filepath.Join("vendor", "bin", "mago"),
		},
	)

	assert.Equal(t, tool.IdMago, mago.Id())
}

func TestMago_DetectMago_System(t *testing.T) {
	mago := test.RunDetectToolLastFound(
		t,
		magoDetect,
		[]string{
			filepath.Join("vendor", "bin", "mago"),
			"mago",
		},
	)

	assert.Equal(t, tool.IdMago, mago.Id())
}

func TestMago_DetectMago_NotFound(t *testing.T) {
	mago := test.RunDetectToolNotFound(
		t,
		magoDetect,
		[]string{
			filepath.Join("vendor", "bin", "mago"),
			"mago",
		},
	)

	assert.Equal(t, tool.IdMago, mago.Id())
}

func TestMago_DetectMago_Config(t *testing.T) {
	mago := test.RunDetectToolConfig(t, magoDetect, "mago")

	assert.Equal(t, tool.IdMago, mago.Id())
}

func TestMago_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		magoLinter,
		internal.ProjectLinterOptions{},
		[]string{"lint"},
	)
}

func TestMago_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		magoLinter,
		internal.ProjectLinterOptions{},
		[]string{"lint"},
	)
}

func TestMago_LintFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesSuccess(
		t,
		magoLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"lint"}, fileNames...),
	)
}

func TestMago_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesFailed(
		t,
		magoLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"lint"}, fileNames...),
	)
}

func TestMago_AnalyzeFiles_All(t *testing.T) {
	test.InvokeAnalyzeFilesSuccess(
		t,
		magoAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"analyze"},
	)
}

func TestMago_AnalyzeFiles_All_Failed(t *testing.T) {
	test.InvokeAnalyzeFilesFailed(
		t,
		magoAnalyzer,
		internal.ProjectAnalyzerOptions{},
		[]string{"analyze"},
	)
}

func TestMago_AnalyzeFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesSuccess(
		t,
		magoAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		append([]string{"analyze"}, fileNames...),
	)
}

func TestMago_AnalyzeFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeAnalyzeFilesFailed(
		t,
		magoAnalyzer,
		internal.ProjectAnalyzerOptions{Filenames: &fileNames},
		append([]string{"analyze"}, fileNames...),
	)
}

func TestMago_AnalyzeFiles_Reporter_NoErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		magoAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyze", "--reporting-format=sarif"},
		"{}",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestMago_AnalyzeFiles_Reporter_SomeErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		magoAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyze", "--reporting-format=sarif"},
		`{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": [
		  {
			"tool": {
			  "driver": {
				"informationUri": "https://github.com/carthage-software/mago",
				"name": "mago",
				"version": "1.57.0"
			  }
			},
			"artifacts": [
			  {
				"location": {
				  "uri": "file:///test.php"
				}
			  }
			],
			"results": [
			  {
				"level": "error",
				"message": {
				  "text": "'x' is assigned a value but never used."
				},
				"locations": [
				  {
					"physicalLocation": {
					  "artifactLocation": {
						"uri": "file:///test.php",
						"index": 0
					  },
					  "region": {
						"startLine": 1,
						"startColumn": 5
					  }
					}
				  }
				],
				"ruleId": "no-unused-vars",
				"ruleIndex": 0
			  }
			]
		  }
		]
  		}`,
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.NotEmpty(t, result.String())
}

func TestMago_AnalyzeFiles_Reporter_NoOutput(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		magoAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyze", "--reporting-format=sarif"},
		"",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestMago_AnalyzeFiles_Reporter_InvalidJson(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeAnalyzeFilesReporter(
		t,
		magoAnalyzer,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"analyze", "--reporting-format=sarif"},
		".",
		func(err error) {
			require.ErrorContains(t, err, "failed to parse SARIF output:")
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestMago_AnalyzeFiles_OutputFile(t *testing.T) {
	test.InvokeAnalyzeFilesOutputFile(
		t,
		magoAnalyzer,
		".",
		[]string{"analyze"},
		nil,
	)
}

func TestMago_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		magoFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"format"},
	)
}

func TestMago_FormatFiles(t *testing.T) {
	fileNames := []string{"tests/*"}

	test.InvokeFormatFilesSuccess(
		t,
		magoFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Filenames:   &fileNames,
		},
		append([]string{"format"}, fileNames...),
	)
}

func TestMago_FormatFiles_CheckAll(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		magoFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
		},
		[]string{"format", "--dry-run"},
	)
}

func TestMago_FormatFiles_Check(t *testing.T) {
	fileNames := []string{"tests/*", "/path/to/file.php"}

	test.InvokeFormatFilesSuccess(
		t,
		magoFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
			Filenames:   &fileNames,
		},
		append([]string{"format", "--dry-run"}, fileNames...),
	)
}

func magoDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectMago(ctx, options)
}

func magoLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewMago(executable)
}

func magoAnalyzer(executable func() (*internal.Executable, error)) internal.ProjectAnalyzer {
	return tool.NewMago(executable)
}

func magoFormatter(executable func() (*internal.Executable, error)) internal.ProjectFormatter {
	return tool.NewMago(executable)
}
