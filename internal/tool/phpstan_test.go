package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPStan_DetectPHPStan_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/phpstan").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPStan(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpstan", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPStan_DetectPHPStan_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/phpstan").
		Return(nil)

	// System installation
	env.OnFindExecutable("phpstan").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPStan(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpstan", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPStan_DetectPHPStan_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/phpstan").
		Return(nil)
	env.OnFindExecutable("phpstan").
		Return(nil)

	tool := DetectPHPStan(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpstan", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse"}).
		Return(nil)

	err := tool.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPStan(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"analyse", "file.php", "/path/to/file2.php"}).
		Return(nil)

	err := tool.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"file.php", "/path/to/file2.php"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
