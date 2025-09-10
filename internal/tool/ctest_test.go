package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCTest_DetectCTest(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("ctest").
		Return(env.NewExecutable("/bin/ctest"))

	tool := DetectCTest(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "ctest", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/ctest", executable.Path)
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
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestCTest_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewCTest(exec.LazyExecutable("ctest"))

	desc := internal.ProjectInfo{Directory: "project", IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(nil)

	err := tool.RunForProject(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestCTest_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewCTest(exec.LazyExecutable("ctest"))

	desc := internal.ProjectInfo{Directory: "project", IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(errors.New("failed"))

	err := tool.RunForProject(desc, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
