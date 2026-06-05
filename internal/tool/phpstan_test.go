package tool_test

import (
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPStan_DetectPHPStan_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpstan")).
		Return(env.NewExecutable("/bin/phpstan"), nil)

	phpStan := tool.DetectPHPStan(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpStan)
	assert.Equal(t, "phpstan", phpStan.Id())
	assert.True(t, phpStan.IsAvailable())

	if executable := phpStan.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpstan", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPStan_DetectPHPStan_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpstan")).
		Return(nil, nil)

	// System installation
	env.OnFindExecutable("phpstan").
		Return(env.NewExecutable("/bin/phpstan"), nil)

	phpStan := tool.DetectPHPStan(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpStan)
	assert.Equal(t, "phpstan", phpStan.Id())
	assert.True(t, phpStan.IsAvailable())

	if executable := phpStan.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpstan", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPStan_DetectPHPStan_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpstan")).
		Return(nil, nil)
	env.OnFindExecutable("phpstan").
		Return(nil, nil)

	phpStan := tool.DetectPHPStan(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpStan)
	assert.Equal(t, "phpstan", phpStan.Id())
	assert.False(t, phpStan.IsAvailable())
	assert.Nil(t, phpStan.Executable())

	env.AssertExpectations(t)
}

func TestPHPStan_DetectPHPStan_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("phpstan-2").
		Return(env.NewExecutable("/bin/phpstan"), nil)

	phpStan := tool.DetectPHPStan(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"phpstan": "phpstan-2"},
	})
	assert.NotNil(t, phpStan)
	assert.Equal(t, "phpstan", phpStan.Id())
	assert.True(t, phpStan.IsAvailable())

	if executable := phpStan.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpstan", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse"}).
		Return(nil)

	err := phpStan.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "file.php", "/path/to/file2.php"}).
		Return(nil)

	err := phpStan.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file.php", "/path/to/file2.php"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat_Json(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "--error-format=json"}).
		Return(nil)

	err := phpStan.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatJson,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat_JUnit(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "--error-format=junit"}).
		Return(nil)

	err := phpStan.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatJUnit,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat_GitHub(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "--error-format=github"}).
		Return(nil)

	err := phpStan.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatGitHub,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat_GitLab(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "--error-format=gitlab"}).
		Return(nil)

	err := phpStan.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatGitLab,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat_TeamCity(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "--error-format=teamcity"}).
		Return(nil)

	err := phpStan.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatTeamCity,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []struct {
		Format internal.LintReportFormat
	}{
		{"test-unsupported"},
	}

	for _, data := range dataSet {
		t.Run(string(data.Format), func(t *testing.T) {
			exec := test.NewExecutable(t)

			phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

			info := internal.ProjectInfo{Directory: "."}

			err := phpStan.LintFiles(
				t.Context(),
				internal.ProjectLinterOptions{
					ProjectInfo: info,
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
