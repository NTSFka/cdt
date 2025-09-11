package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParaTest_DetectParaTest_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/paratest").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectParaTest(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "paratest", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/paratest").
		Return(nil)

	// System installation
	env.OnFindExecutable("paratest").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectParaTest(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "paratest", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/paratest").
		Return(nil)
	env.OnFindExecutable("paratest").
		Return(nil)

	tool := DetectParaTest(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "paratest", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestParaTest_ParaTest_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := tool.TestAll(context.Background(), internal.ProjectTesterOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_ParaTest_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := tool.TestPattern(context.Background(), internal.ProjectTesterOptions{ProjectInfo: info}, "tests/*")
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
