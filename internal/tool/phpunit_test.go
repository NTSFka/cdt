package tool_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPUnit_DetectPHPUnit_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/phpunit").
		Return(env.NewExecutable("/bin/phpunit"))

	phpUnit := tool.DetectPHPUnit(t.Context(), env)
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
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
		Return(env.NewExecutable("/bin/phpunit"))

	phpUnit := tool.DetectPHPUnit(t.Context(), env)
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/phpunit").
		Return(nil)
	env.OnFindExecutable("phpunit").
		Return(nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), env)
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.False(t, phpUnit.IsAvailable())
	assert.Nil(t, phpUnit.Executable())

	env.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := phpUnit.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := phpUnit.TestPattern(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info}, "tests/*")
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
