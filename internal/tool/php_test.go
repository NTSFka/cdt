package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHP_DetectPHP(t *testing.T) {
	php := test.RunDetectToolLastFound(t, phpDetect, []string{"php"})

	assert.Equal(t, tool.IdPHP, php.Id())
}

func TestPHP_DetectPHP_NotFound(t *testing.T) {
	php := test.RunDetectToolNotFound(t, phpDetect, []string{"php"})

	assert.Equal(t, tool.IdPHP, php.Id())
}

func TestPHP_DetectPHP_Config(t *testing.T) {
	php := test.RunDetectToolConfig(t, phpDetect, "php")

	assert.Equal(t, tool.IdPHP, php.Id())
}

func TestPHP_PHP_RunTarget_Success(t *testing.T) {
	test.InvokeRunTargetSuccess(
		t,
		phpRunner,
		internal.ProjectRunnerOptions{},
		"index.php",
		[]string{"-f", "index.php"},
	)
}

func TestPHP_PHP_RunTarget_Failed(t *testing.T) {
	test.InvokeRunTargetFailed(
		t,
		phpRunner,
		internal.ProjectRunnerOptions{},
		"index.php",
		[]string{"-f", "index.php"},
	)
}

func TestPHP_PHP_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesError(
		t,
		phpLinter,
		internal.ProjectLinterOptions{},
		func(err error) {
			require.EqualError(t, err, "no files to lint")
		},
	)
}

func TestPHP_PHP_LintFiles(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesSuccess(
		t,
		phpLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"-l"}, fileNames...),
	)
}

func TestPHP_PHP_LintFiles_Failed(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesFailed(
		t,
		phpLinter,
		internal.ProjectLinterOptions{Filenames: &fileNames},
		append([]string{"-l"}, fileNames...),
	)
}

func TestPHP_PHP_LintFiles_OutputFormat_Raw(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesOutputFormatCheck(
		t,
		phpLinter,
		".",
		internal.LintReportFormatRaw,
		append([]string{"-l"}, fileNames...),
		&fileNames,
	)
}

func TestPHP_PHP_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.InvokeLintFilesOutputFormatUnsupported(
				t,
				phpLinter,
				".",
				format,
				&[]string{"file.php", "/path/to/file2.php"},
			)
		})
	}
}

func TestPHP_PHP_LintFiles_OutputFile(t *testing.T) {
	fileNames := []string{"file.php", "/path/to/file2.php"}

	test.InvokeLintFilesOutputFile(
		t,
		phpLinter,
		".",
		append([]string{"-l"}, fileNames...),
		&fileNames,
	)
}

func phpDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHP(ctx, options)
}

func phpRunner(executable func() (*internal.Executable, error)) internal.ProjectRunner {
	return tool.NewPHP(executable)
}

func phpLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewPHP(executable)
}
