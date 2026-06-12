package tool_test

import (
	"context"
	"errors"
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
	test.RunDetectToolNotFound(t, phpDetect, []string{"php"})
}

func TestPHP_DetectPHP_Config(t *testing.T) {
	php := test.RunDetectToolConfig(t, phpDetect, "php")

	assert.Equal(t, tool.IdPHP, php.Id())
}

func TestPHP_PHP_RunTarget(t *testing.T) {
	exec := test.NewExecutable(t)

	php := tool.NewPHP(exec.LazyExecutable("php"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("php", []string{"-f", "index.php"}).
		Return(nil)

	err := php.RunTarget(t.Context(), internal.ProjectRunnerOptions{ProjectInfo: info}, "index.php")
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHP_PHP_RunTarget_Fail(t *testing.T) {
	exec := test.NewExecutable(t)

	php := tool.NewPHP(exec.LazyExecutable("php"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("php", []string{"-f", "index.php"}).
		Return(errors.New("failed"))

	err := php.RunTarget(t.Context(), internal.ProjectRunnerOptions{ProjectInfo: info}, "index.php")
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestPHP_PHP_LintFiles_All(t *testing.T) {
	test.RunLintFilesInvoked(
		t,
		phpBuildLinter,
		internal.ProjectLinterOptions{},
		nil,
		func(err error) {
			require.EqualError(t, err, "no files to lint")
		},
	)
}

func TestPHP_PHP_LintFiles(t *testing.T) {
	test.RunLintFilesSuccess(
		t,
		phpBuildLinter,
		internal.ProjectLinterOptions{
			Filenames: &[]string{"file.php", "/path/to/file2.php"},
		},
		[]string{"-l", "file.php", "/path/to/file2.php"},
	)
}

func TestPHP_PHP_LintFiles_OutputFormat_Raw(t *testing.T) {
	dataSet := []struct {
		format internal.LintReportFormat
	}{
		{internal.LintReportFormatRaw},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			test.RunLintFilesOutputFormatCheck(
				t,
				phpBuildLinter,
				data.format,
				[]string{"-l", "file.php", "/path/to/file2.php"},
				&[]string{"file.php", "/path/to/file2.php"},
			)
		})
	}
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
			test.RunLintFilesOutputFormatUnsupported(
				t,
				phpBuildLinter,
				format,
				&[]string{"file.php", "/path/to/file2.php"},
			)
		})
	}
}

func TestPHP_PHP_LintFiles_OutputFile(t *testing.T) {
	test.RunLintFilesOutputFile(
		t,
		phpBuildLinter,
		[]string{"-l", "file.php", "/path/to/file2.php"},
		&[]string{"file.php", "/path/to/file2.php"},
	)
}

func phpDetect(ctx context.Context, options tool.DetectOptions) internal.Tool {
	return tool.DetectPHP(ctx, options)
}

func phpBuildLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewPHP(executable)
}
