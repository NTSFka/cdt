package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPHPUnit_DetectPHPUnit_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/phpunit").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPUnit(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpunit", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/phpunit").
		Return(nil)

	// System installation
	env.OnFindExecutable("phpunit").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPUnit(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpunit", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/phpunit").
		Return(nil)
	env.OnFindExecutable("phpunit").
		Return(nil)

	tool := DetectPHPUnit(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpunit", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPUnit(exec.LazyExecutable("test"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := tool.TestAll(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPUnit(exec.LazyExecutable("test"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := tool.Test(desc, "tests/*", []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
