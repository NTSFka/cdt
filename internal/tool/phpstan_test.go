package tool_test

import (
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
	env.OnFindExecutable("vendor/bin/phpstan").
		Return(env.NewExecutable("/bin/phpstan"), nil)

	phpStan := tool.DetectPHPStan(t.Context(), env)
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
	env.OnFindExecutable("vendor/bin/phpstan").
		Return(nil, nil)

	// System installation
	env.OnFindExecutable("phpstan").
		Return(env.NewExecutable("/bin/phpstan"), nil)

	phpStan := tool.DetectPHPStan(t.Context(), env)
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

	env.OnFindExecutable("vendor/bin/phpstan").
		Return(nil, nil)
	env.OnFindExecutable("phpstan").
		Return(nil, nil)

	phpStan := tool.DetectPHPStan(t.Context(), env)
	assert.NotNil(t, phpStan)
	assert.Equal(t, "phpstan", phpStan.Id())
	assert.False(t, phpStan.IsAvailable())
	assert.Nil(t, phpStan.Executable())

	env.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse"}).
		Return(nil)

	err := phpStan.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	phpStan := tool.NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "file.php", "/path/to/file2.php"}).
		Return(nil)

	err := phpStan.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"file.php", "/path/to/file2.php"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
