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
	test.InvokeLintFilesSuccess(
		t,
		myPyLinter,
		internal.ProjectLinterOptions{},
		[]string{"*.py"},
	)
}

func TestMyPy_MyPy_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		myPyLinter,
		internal.ProjectLinterOptions{},
		[]string{"*.py"},
	)
}

func TestMyPy_MyPy_LintFiles(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesSuccess(
		t,
		myPyLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestMyPy_MyPy_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.py", "/path/to/file2.py"}

	test.InvokeLintFilesFailed(
		t,
		myPyLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestMyPy_MyPy_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		myPyLinter,
		internal.LintReportFormatRaw,
		[]string{"*.py"},
		nil,
	)
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
			test.InvokeLintFilesOutputFormatUnsupported(t, myPyLinter, format, nil)
		})
	}
}

func TestMyPy_MyPy_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(t, myPyLinter, []string{"*.py"}, nil)
}

func myPyDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectMyPy(ctx, options)
}

func myPyLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewMyPy(executable)
}
