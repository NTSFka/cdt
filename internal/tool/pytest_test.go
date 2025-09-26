package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPyTest_DetectPyTest(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPyTest(t.Context(), env)
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

	tool := DetectPyTest(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pytest", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPyTest_PyTest_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := tool.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPyTest_PyTest_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := tool.TestPattern(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info}, "tests/*")
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
