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

func TestPHPParallelLint_DetectPHPParallelLint_Composer(t *testing.T) {
	phpParallelLint := test.RunDetectToolLastFound(
		t,
		phpParallelLintDetect,
		[]string{
			filepath.Join("vendor", "bin", "parallel-lint"),
		},
	)

	assert.Equal(t, tool.IdPHPParallelLint, phpParallelLint.Id())
}

func TestPHPParallelLint_DetectPHPParallelLint_System(t *testing.T) {
	phpParallelLint := test.RunDetectToolLastFound(
		t,
		phpParallelLintDetect,
		[]string{
			filepath.Join("vendor", "bin", "parallel-lint"),
			"parallel-lint",
		},
	)

	assert.Equal(t, tool.IdPHPParallelLint, phpParallelLint.Id())
}

func TestPHPParallelLint_DetectPHPParallelLint_NotFound(t *testing.T) {
	phpParallelLint := test.RunDetectToolNotFound(
		t,
		phpParallelLintDetect,
		[]string{
			filepath.Join("vendor", "bin", "parallel-lint"),
			"parallel-lint",
		},
	)

	assert.Equal(t, tool.IdPHPParallelLint, phpParallelLint.Id())
}

func TestPHPParallelLint_DetectPHPParallelLint_Config(t *testing.T) {
	phpParallelLint := test.RunDetectToolConfig(t, phpParallelLintDetect, "php-parallel-lint")

	assert.Equal(t, tool.IdPHPParallelLint, phpParallelLint.Id())
}

func TestPHPParallelLint_PHPParallelLint_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		phpParallelLintLinter,
		internal.ProjectLinterOptions{},
		[]string{},
	)
}

func TestPHPParallelLint_PHPParallelLint_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		phpParallelLintLinter,
		internal.ProjectLinterOptions{},
		[]string{},
	)
}

func TestPHPParallelLint_PHPParallelLint_LintFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesSuccess(
		t,
		phpParallelLintLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func TestPHPParallelLint_PHPParallelLint_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesFailed(
		t,
		phpParallelLintLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		fileNames,
	)
}

func phpParallelLintDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHPParallelLint(ctx, options)
}

func phpParallelLintLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewPHPParallelLint(executable)
}
