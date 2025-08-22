package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPyTest_DetectPyTest(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPyTest(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pytest", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}

	env.AssertExpectations(t)
}

func TestPyTest_DetectPyTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(nil)

	tool := DetectPyTest(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pytest", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestPyTest_PyTest_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPyTest(exec.LazyExecutable("test"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Tester: tool})

	exec.OnRun("test", []string{}).
		Return(nil)

	err := tool.TestAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPyTest_PyTest_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPyTest(exec.LazyExecutable("test"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Tester: tool})

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := tool.Test(p, "tests/*", []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
