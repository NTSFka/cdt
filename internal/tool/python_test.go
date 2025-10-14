package tool_test

import (
	"errors"
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
	env := test.NewEnvironment(t)

	env.OnFindExecutable("python3").
		Return(env.NewExecutable("/bin/python"), nil)

	python := tool.DetectPython(t.Context(), env)
	assert.NotNil(t, python)
	assert.Equal(t, "python", python.Id())
	assert.True(t, python.IsAvailable())

	if executable := python.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/python", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPython_DetectPython_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("python3").
		Return(nil, nil)

	python := tool.DetectPython(t.Context(), env)
	assert.NotNil(t, python)
	assert.Equal(t, "python", python.Id())
	assert.False(t, python.IsAvailable())
	assert.Nil(t, python.Executable())

	env.AssertExpectations(t)
}

func TestPython_Python_RunTarget(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python3"))

	info := internal.ProjectInfo{Directory: "."}

	executable.OnRun("python3", []string{"main.py"}).
		Return(nil)

	err := python.RunTarget(
		t.Context(),
		internal.ProjectRunnerOptions{ProjectInfo: info},
		"main.py",
	)
	require.NoError(t, err)

	executable.AssertExpectations(t)
}

func TestPython_Python_RunTarget_Fail(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python3"))

	info := internal.ProjectInfo{Directory: "."}

	executable.OnRun("python3", []string{"main.py"}).
		Return(errors.New("failed"))

	err := python.RunTarget(
		t.Context(),
		internal.ProjectRunnerOptions{ProjectInfo: info},
		"main.py",
	)
	require.EqualError(t, err, "failed")

	executable.AssertExpectations(t)
}

func TestPython_CreateEnvironment(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python"))
	assert.NotNil(t, python)
	assert.Equal(t, []string{"pyenv"}, python.Aliases())

	env, err := python.CreateEnvironment(".", ".venv")
	require.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "pyenv", env.Id())

	executable.AssertExpectations(t)
}

func TestPython_CreateEnvironment_NoPath(t *testing.T) {
	executable := test.NewExecutable(t)

	python := tool.NewPython(executable.LazyExecutable("python"))
	assert.NotNil(t, python)
	assert.Equal(t, []string{"pyenv"}, python.Aliases())

	env, err := python.CreateEnvironment(".", "")
	require.EqualError(t, err, "python virtual environment path is required")
	assert.Nil(t, env)

	executable.AssertExpectations(t)
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
	_, err = os.OpenFile(filepath.Join(testDir, "bin", "tool1"), os.O_RDONLY|os.O_CREATE, 0600)
	require.NoError(t, err)

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
	_, err = os.OpenFile(filepath.Join(testDir, "Scripts", "tool1"), os.O_RDONLY|os.O_CREATE, 0600)
	require.NoError(t, err)

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
