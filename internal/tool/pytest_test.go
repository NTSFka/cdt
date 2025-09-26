package tool_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPyTest_DetectPyTest(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(env.NewExecutable("/bin/pytest"))

	pyTest := tool.DetectPyTest(t.Context(), env)
	assert.NotNil(t, pyTest)
	assert.Equal(t, "pytest", pyTest.Id())
	assert.True(t, pyTest.IsAvailable())

	if executable := pyTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pytest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPyTest_DetectPyTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(nil)

	pyTest := tool.DetectPyTest(t.Context(), env)
	assert.NotNil(t, pyTest)
	assert.Equal(t, "pytest", pyTest.Id())
	assert.False(t, pyTest.IsAvailable())
	assert.Nil(t, pyTest.Executable())

	env.AssertExpectations(t)
}

func TestPyTest_PyTest_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	pyTest := tool.NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := pyTest.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPyTest_PyTest_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	pyTest := tool.NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := pyTest.TestPattern(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info}, "tests/*")
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
