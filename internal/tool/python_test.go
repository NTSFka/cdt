package tool_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPython_DetectPython(t *testing.T) {
	python := test.RunDetectToolLastFound(t, pythonDetect, []string{"python3"})

	assert.Equal(t, tool.IdPython, python.Id())
}

func TestPython_DetectPython_NotFound(t *testing.T) {
	python := test.RunDetectToolNotFound(t, pythonDetect, []string{"python3"})

	assert.Equal(t, tool.IdPython, python.Id())
}

func TestPython_DetectPython_Config(t *testing.T) {
	python := test.RunDetectToolConfig(t, pythonDetect, "python")

	assert.Equal(t, tool.IdPython, python.Id())
}

func TestPython_Python_RunTarget(t *testing.T) {
	test.InvokeRunTargetSuccess(
		t,
		pythonRunner,
		internal.ProjectRunnerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		"main.py",
		[]string{"main.py"},
	)
}

func TestPython_Python_RunTarget_Failed(t *testing.T) {
	test.InvokeRunTargetFailed(
		t,
		pythonRunner,
		internal.ProjectRunnerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		"main.py",
		[]string{"main.py"},
	)
}

func TestPython_CreateEnvironment_Success(t *testing.T) {
	python := test.InvokeCreateEnvironment(
		t,
		pythonEnvironmentProvider,
		".",
		".venv",
		func(exec *test.Executable, executablePath string) {
			// No invocation
		},
		func(env internal.Environment, err error) {
			require.NoError(t, err)
			require.NotNil(t, env)
			assert.Equal(t, "pyenv", env.Id())
		},
	)

	assert.Equal(t, []string{"pyenv"}, python.Aliases())
}

func TestPython_CreateEnvironment_NoPath(t *testing.T) {
	test.InvokeCreateEnvironment(
		t,
		pythonEnvironmentProvider,
		".",
		"",
		func(exec *test.Executable, executablePath string) {
			// No invocation
		},
		func(env internal.Environment, err error) {
			require.EqualError(t, err, "python virtual environment path is required")
			require.Nil(t, env)
		},
	)
}

func TestPython_Environment_Start(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python"))
	assert.NotNil(t, python)

	env, err := python.CreateEnvironment(".", ".venv")
	require.NoError(t, err)

	executable.OnRun("python", []string{"-m", "venv", ".venv"}).Return(nil)

	err = env.Start(t.Context())
	require.NoError(t, err)

	executable.AssertExpectations(t)
}

func TestPython_Environment_IsRunning_True(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python"))
	assert.NotNil(t, python)

	env, err := python.CreateEnvironment(".", ".venv")
	require.NoError(t, err)

	result := env.IsRunning(t.Context())
	assert.True(t, result)

	executable.AssertExpectations(t)
}

func TestPython_Environment_Stop(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python"))
	assert.NotNil(t, python)

	env, err := python.CreateEnvironment(".", ".venv")
	require.NoError(t, err)

	err = env.Stop(t.Context())
	require.NoError(t, err)

	executable.AssertExpectations(t)
}

func TestPython_Environment_Cleanup(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python"))
	assert.NotNil(t, python)

	env, err := python.CreateEnvironment(".", ".venv")
	require.NoError(t, err)

	err = env.Cleanup(t.Context())
	require.NoError(t, err)

	executable.AssertExpectations(t)
}

func TestPython_Environment_FindExecutable(t *testing.T) {
	python := tool.NewPython(test.LazyExecutable("python"))
	assert.NotNil(t, python)

	testDir := t.TempDir()

	env, err := python.CreateEnvironment(".", testDir)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(testDir, "bin"), 0700)
	require.NoError(t, err)
	file, err := os.Create(filepath.Join(testDir, "bin", "tool1"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	executable, err := env.FindExecutable(t.Context(), "tool1")
	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.Equal(t, "tool1", executable.Path)
}

func TestPython_Environment_FindExecutable_Windows(t *testing.T) {
	python := tool.NewPython(test.LazyExecutable("python"))
	assert.NotNil(t, python)

	testDir := t.TempDir()

	env, err := python.CreateEnvironment(".", testDir)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(testDir, "Scripts"), 0700)
	require.NoError(t, err)
	file, err := os.Create(filepath.Join(testDir, "Scripts", "tool1"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	executable, err := env.FindExecutable(t.Context(), "tool1")
	require.NotNil(t, executable)
	require.NoError(t, err)
	assert.Equal(t, "tool1", executable.Path)
}

func TestPython_Environment_FindExecutable_Failed(t *testing.T) {
	python := tool.NewPython(test.LazyExecutable("python"))
	assert.NotNil(t, python)

	testDir := t.TempDir()

	env, err := python.CreateEnvironment(".", testDir)
	require.NoError(t, err)

	executable, err := env.FindExecutable(t.Context(), "tool1")
	assert.Nil(t, executable)
	assert.NoError(t, err)
}

func TestPython_Environment_RunExecutable_NotFound(t *testing.T) {
	python := tool.NewPython(test.LazyExecutable("python"))
	assert.NotNil(t, python)

	testDir := t.TempDir()

	env, err := python.CreateEnvironment(".", testDir)
	require.NoError(t, err)

	err = env.RunExecutable(t.Context(), internal.RunOptions{}, "echo", []string{"arg1", "arg2"})
	require.EqualError(t, err, "executable not found")
}

func pythonDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPython(ctx, options)
}

func pythonRunner(executable func() (*internal.Executable, error)) internal.ProjectRunner {
	return tool.NewPython(executable)
}

func pythonEnvironmentProvider(
	executable func() (*internal.Executable, error),
) internal.EnvironmentProvider {
	return tool.NewPython(executable)
}
