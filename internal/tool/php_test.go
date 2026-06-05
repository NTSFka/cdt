package tool_test

import (
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHP_DetectPHP(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php").
		Return(env.NewExecutable("/bin/php"), nil)

	php := tool.DetectPHP(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, php)
	assert.Equal(t, "php", php.Id())
	assert.True(t, php.IsAvailable())

	if executable := php.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/php", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHP_DetectPHP_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php").
		Return(nil, nil)

	php := tool.DetectPHP(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, php)
	assert.Equal(t, "php", php.Id())
	assert.False(t, php.IsAvailable())
	assert.Nil(t, php.Executable())

	env.AssertExpectations(t)
}

func TestPHP_DetectPHP_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php-8.2").
		Return(env.NewExecutable("/bin/php"), nil)

	php := tool.DetectPHP(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"php": "php-8.2"},
	})
	assert.NotNil(t, php)
	assert.Equal(t, "php", php.Id())
	assert.True(t, php.IsAvailable())

	if executable := php.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/php", executable.Path)
	}

	env.AssertExpectations(t)
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
	exec := test.NewExecutable(t)

	linter := tool.NewPHP(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	err := linter.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.EqualError(t, err, "no files to lint")

	exec.AssertExpectations(t)
}

func TestPHP_PHP_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	linter := tool.NewPHP(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"-l", "file.php", "/path/to/file2.php"}).
		Return(nil)

	err := linter.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file.php", "/path/to/file2.php"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHP_PHP_LintFiles_OutputFormat_Raw(t *testing.T) {
	exec := test.NewExecutable(t)

	linter := tool.NewPHP(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"-l", "file.php", "/path/to/file2.php"}).
		Return(nil)

	err := linter.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file.php", "/path/to/file2.php"},
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatRaw,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHP_PHP_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []struct {
		Format internal.LintReportFormat
	}{
		{internal.LintReportFormatJson},
		{internal.LintReportFormatJUnit},
		{internal.LintReportFormatGitHub},
		{internal.LintReportFormatGitLab},
		{internal.LintReportFormatTeamCity},
		{"test-unsupported"},
	}

	for _, data := range dataSet {
		t.Run(string(data.Format), func(t *testing.T) {
			exec := test.NewExecutable(t)

			linter := tool.NewPHP(exec.LazyExecutable("lint"))

			info := internal.ProjectInfo{Directory: "."}

			err := linter.LintFiles(
				t.Context(),
				internal.ProjectLinterOptions{
					ProjectInfo: info,
					Filenames:   &[]string{"file.php", "/path/to/file2.php"},
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: data.Format,
					},
				},
			)
			require.EqualError(t, err, "unsupported report format: "+string(data.Format))

			exec.AssertExpectations(t)
		})
	}
}
