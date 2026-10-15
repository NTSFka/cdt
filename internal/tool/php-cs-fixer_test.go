package tool_test

import (
	"context"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestPHPCSFixer_DetectPHPCSFixer_Composer(t *testing.T) {
	phpcsFixer := test.RunDetectToolLastFound(
		t,
		phpcsFixerDetect,
		[]string{filepath.Join("vendor", "bin", "php-cs-fixer")},
	)

	assert.Equal(t, tool.IdPHPCSFixer, phpcsFixer.Id())
}

func TestPHPCSFixer_DetectPHPCSFixer_System(t *testing.T) {
	phpcsFixer := test.RunDetectToolLastFound(
		t,
		phpcsFixerDetect,
		[]string{filepath.Join("vendor", "bin", "php-cs-fixer"), "php-cs-fixer"},
	)

	assert.Equal(t, tool.IdPHPCSFixer, phpcsFixer.Id())
}

func TestPHPCSFixer_DetectPHPCSFixer_NotFound(t *testing.T) {
	phpcsFixer := test.RunDetectToolNotFound(
		t,
		phpcsFixerDetect,
		[]string{filepath.Join("vendor", "bin", "php-cs-fixer"), "php-cs-fixer"},
	)

	assert.Equal(t, tool.IdPHPCSFixer, phpcsFixer.Id())
}

func TestPHPCSFixer_DetectPHPCSFixer_Config(t *testing.T) {
	phpcsFixer := test.RunDetectToolConfig(t, phpcsFixerDetect, "php-cs-fixer")

	assert.Equal(t, tool.IdPHPCSFixer, phpcsFixer.Id())
}

func TestPHPCSFixer_PHPCSFixer_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		phpcsFixerFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"fix"},
	)
}

func TestPHPCSFixer_PHPCSFixer_FormatFiles(t *testing.T) {
	fileNames := []string{"tests/*"}

	test.InvokeFormatFilesSuccess(
		t,
		phpcsFixerFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Filenames:   &fileNames,
		},
		append([]string{"fix"}, fileNames...),
	)
}

func TestPHPCSFixer_PHPCSFixer_FormatFiles_CheckAll(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		phpcsFixerFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
		},
		[]string{"fix", "--dry-run"},
	)
}

func TestPHPCSFixer_PHPCSFixer_FormatFiles_Check(t *testing.T) {
	fileNames := []string{"tests/*", "/path/to/file.php"}

	test.InvokeFormatFilesSuccess(
		t,
		phpcsFixerFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
			Filenames:   &fileNames,
		},
		append([]string{"fix", "--dry-run"}, fileNames...),
	)
}

func phpcsFixerDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHPCSFixer(ctx, options)
}

func phpcsFixerFormatter(
	executable func() (*internal.Executable, error),
) internal.ProjectFormatter {
	return tool.NewPHPCSFixer(executable)
}
