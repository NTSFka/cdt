package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPyTest_DetectPyTest(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPyTest(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pytest", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPyTest_DetectPyTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(nil)

	tool := DetectPyTest(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pytest", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPyTest_PyTest_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPyTest(exec.LazyExecutable("test"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := tool.TestAll(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPyTest_PyTest_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPyTest(exec.LazyExecutable("test"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := tool.Test(context.Background(), desc, "tests/*", []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
