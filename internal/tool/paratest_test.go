package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParaTest_DetectParaTest_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/paratest").
		Return(env.NewExecutable("/bin/paratest"), nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.True(t, paraTest.IsAvailable())

	if executable := paraTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/paratest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("vendor/bin/paratest").
		Return(nil, nil)

	// System installation
	env.OnFindExecutable("paratest").
		Return(env.NewExecutable("/bin/paratest"), nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.True(t, paraTest.IsAvailable())

	if executable := paraTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/paratest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("vendor/bin/paratest").
		Return(nil, nil)
	env.OnFindExecutable("paratest").
		Return(nil, nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.False(t, paraTest.IsAvailable())
	assert.Nil(t, paraTest.Executable())

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("paratest-7").
		Return(env.NewExecutable("/bin/paratest"), nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"paratest": "paratest-7"},
	})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.True(t, paraTest.IsAvailable())

	if executable := paraTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/paratest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_ParaTest_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := paraTest.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_ParaTest_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := paraTest.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: info},
		"tests/*",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
