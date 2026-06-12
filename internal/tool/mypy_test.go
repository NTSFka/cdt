package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestMyPy_DetectMyPy(t *testing.T) {
	myPy := test.RunDetectToolLastFound(t, myPyDetect, []string{"mypy"})

	assert.Equal(t, tool.IdMyPy, myPy.Id())
}

func TestMyPy_DetectMyPy_NotFound(t *testing.T) {
	myPy := test.RunDetectToolNotFound(t, myPyDetect, []string{"mypy"})

	assert.Equal(t, tool.IdMyPy, myPy.Id())
}

func TestMyPy_DetectMyPy_Config(t *testing.T) {
	myPy := test.RunDetectToolConfig(t, myPyDetect, "mypy")

	assert.Equal(t, tool.IdMyPy, myPy.Id())
}

func TestMyPy_MyPy_LintFiles_All(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		myPyBuildLinter,
		internal.ProjectLinterOptions{},
		[]string{"*.py"},
	)
}

func TestMyPy_MyPy_LintFiles(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		myPyBuildLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"file.py", "/path/to/file2.py"}},
		[]string{"file.py", "/path/to/file2.py"},
	)
}

func TestMyPy_MyPy_LintFiles_OutputFormat(t *testing.T) {
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
				myPyBuildLinter,
				data.format,
				append([]string{"*.py"}, data.args...),
				nil,
			)
		})
	}
}

func TestMyPy_MyPy_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.RunLintFilesOutputFormatUnsupported(t, myPyBuildLinter, format, nil)
		})
	}
}

func TestMyPy_MyPy_LintFiles_OutputFile(t *testing.T) {
	test.RunLintFilesOutputFile(t, myPyBuildLinter, []string{"*.py"}, nil)
}

func myPyDetect(ctx context.Context, options tool.DetectOptions) internal.Tool {
	return tool.DetectMyPy(ctx, options)
}

func myPyBuildLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewMyPy(executable)
}
