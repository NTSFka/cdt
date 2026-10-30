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

func TestGolangCILint_LintFiles_Reporter_NoErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeLintFilesReporter(
		t,
		golangCILintLinter,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"run", "--output.sarif.path", "stdout"},
		"{}",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestGolangCILint_LintFiles_Reporter_SomeErrors(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeLintFilesReporter(
		t,
		golangCILintLinter,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"run", "--output.sarif.path", "stdout"},
		`{
		"version": "2.1.0",
		"$schema": "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		"runs": [
		  {
			"artifacts": [
			  {
				"location": {
				  "uri": "file:///test.go"
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
						"uri": "file:///test.go",
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

func TestGolangCILint_LintFiles_Reporter_NoOutput(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeLintFilesReporter(
		t,
		golangCILintLinter,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"run", "--output.sarif.path", "stdout"},
		"",
		func(err error) {
			require.NoError(t, err)
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestGolangCILint_LintFiles_Reporter_InvalidJson(t *testing.T) {
	result := bytes.Buffer{}

	test.InvokeLintFilesReporter(
		t,
		golangCILintLinter,
		".",
		reporter.NewSimpleSarifReporter(&result),
		[]string{"run", "--output.sarif.path", "stdout"},
		".",
		func(err error) {
			require.ErrorContains(t, err, "failed to parse SARIF output:")
		},
	)

	assert.Equal(t, "\nTotal issues: 0\n\n", result.String())
}

func TestGolangCILint_GolangCILint_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(
		t,
		golangCILintLinter,
		".",
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
