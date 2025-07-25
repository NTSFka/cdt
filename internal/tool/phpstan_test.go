package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPHPStan_DetectPHPStan_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/phpstan").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPStan(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpstan", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
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

	tool := DetectPHPStan(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpstan", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}

	env.AssertExpectations(t)
}

func TestPHPStan_DetectPHPStan_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/phpstan").
		Return(nil)
	env.OnFindExecutable("phpstan").
		Return(nil)

	tool := DetectPHPStan(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpstan", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestPHPStan_PHPStan_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPStan(exec.LazyExecutable("lint"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"analyse"}).
		Return(nil)

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPStan_PHPStan_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPStan(exec.LazyExecutable("lint"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"analyse", "file.php", "/path/to/file2.php"}).
		Return(nil)

	err := tool.LintFiles(p, []string{"file.php", "/path/to/file2.php"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
