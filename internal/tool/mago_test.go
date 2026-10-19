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

func TestMago_Mago_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		magoLinter,
		internal.ProjectLinterOptions{},
		[]string{"lint"},
	)
}

func TestMago_Mago_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		magoLinter,
		internal.ProjectLinterOptions{},
		[]string{"lint"},
	)
}

func TestMago_Mago_LintFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesSuccess(
		t,
		magoLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"lint"}, fileNames...),
	)
}

func TestMago_Mago_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesFailed(
		t,
		magoLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"lint"}, fileNames...),
	)
}

func TestMago_Mago_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		magoFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"format"},
	)
}

func TestMago_Mago_FormatFiles(t *testing.T) {
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

func TestMago_Mago_FormatFiles_CheckAll(t *testing.T) {
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

func TestMago_Mago_FormatFiles_Check(t *testing.T) {
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

func magoFormatter(executable func() (*internal.Executable, error)) internal.ProjectFormatter {
	return tool.NewMago(executable)
}
