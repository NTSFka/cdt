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

func TestMyPy_MyPy_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(t, myPyLinter, ".", []string{"*.py"}, nil)
}

func myPyDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectMyPy(ctx, options)
}

func myPyLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewMyPy(executable)
}
