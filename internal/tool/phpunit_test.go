package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPUnit_DetectPHPUnit_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/phpunit").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHPUnit(context.Background(), env)
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

	tool := DetectPHPUnit(context.Background(), env)
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

	tool := DetectPHPUnit(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "phpunit", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := tool.TestAll(context.Background(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := tool.TestPattern(context.Background(), internal.ProjectTesterOptions{ProjectInfo: info}, "tests/*")
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
