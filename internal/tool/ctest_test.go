package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCTest_DetectCTest(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("ctest").
		Return(env.NewExecutable("/bin/ctest"))

	tool := DetectCTest(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "ctest", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/ctest", *path)
	}

	env.AssertExpectations(t)
}

func TestCTest_DetectCTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("ctest").
		Return(nil)

	tool := DetectCTest(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "ctest", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestCTest_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewCTest(exec.LazyExecutable("ctest"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(nil)

	err := tool.RunForProject(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestCTest_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewCTest(exec.LazyExecutable("ctest"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	exec.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(errors.New("failed"))

	err := tool.RunForProject(p, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
