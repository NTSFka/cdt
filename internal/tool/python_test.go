package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPython_DetectPython(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("python3").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPython(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "python", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPython_DetectPython_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("python3").
		Return(nil)

	tool := DetectPython(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "python", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPython_Python_RunTarget(t *testing.T) {
	python := test.NewExecutable(t)

	tool := NewPython(python.LazyExecutable("python3"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Runner: tool})

	python.OnRun("python3", []string{"main.py"}).
		Return(nil)

	err := tool.RunTarget(p, "main.py", []string{})
	assert.NoError(t, err)

	python.AssertExpectations(t)
}

func TestPython_Python_RunTarget_Fail(t *testing.T) {
	python := test.NewExecutable(t)

	tool := NewPython(python.LazyExecutable("python3"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Runner: tool})

	python.OnRun("python3", []string{"main.py"}).
		Return(errors.New("failed"))

	err := tool.RunTarget(p, "main.py", []string{})
	assert.EqualError(t, err, "failed")

	python.AssertExpectations(t)
}

func TestPython_CreateEnvironment(t *testing.T) {
	python := test.NewExecutable(t)

	tool := NewPython(python.LazyExecutable("python"))
	assert.NotNil(t, tool)
	assert.Equal(t, []string{"pyenv"}, tool.Aliases())

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "pyenv", env.Id())

	python.AssertExpectations(t)
}

func TestPython_Environment_Start(t *testing.T) {
	python := test.NewExecutable(t)

	tool := NewPython(python.LazyExecutable("python"))
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	python.OnRun("python", []string{"-m", "venv", ".venv"}).Return(nil)

	err = env.Start(context.Background())
	assert.NoError(t, err)

	python.AssertExpectations(t)
}

func TestPython_Environment_IsRunning_True(t *testing.T) {
	python := test.NewExecutable(t)

	tool := NewPython(python.LazyExecutable("python"))
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	result := env.IsRunning(context.Background())
	assert.True(t, result)

	python.AssertExpectations(t)
}

func TestPython_Environment_Stop(t *testing.T) {
	python := test.NewExecutable(t)

	tool := NewPython(python.LazyExecutable("python"))
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	err = env.Stop(context.Background())
	assert.NoError(t, err)

	python.AssertExpectations(t)
}

func TestPython_Environment_Cleanup(t *testing.T) {
	python := test.NewExecutable(t)

	tool := NewPython(python.LazyExecutable("python"))
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	err = env.Cleanup(context.Background())
	assert.NoError(t, err)

	python.AssertExpectations(t)
}

func TestPython_Environment_FindExecutable(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	testDir := t.TempDir()

	env, err := tool.CreateEnvironment(".", testDir)
	assert.NoError(t, err)

	err = os.MkdirAll(filepath.Join(testDir, "bin"), 0755)
	assert.NoError(t, err)
	_, err = os.OpenFile(filepath.Join(testDir, "bin", "tool1"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	executable := env.FindExecutable("tool1")
	if assert.NotNil(t, executable) {
		assert.Equal(t, "tool1", executable.Path)
	}
}

func TestPython_Environment_FindExecutable_Windows(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	testDir := t.TempDir()

	env, err := tool.CreateEnvironment(".", testDir)
	assert.NoError(t, err)

	err = os.MkdirAll(filepath.Join(testDir, "Scripts"), 0755)
	assert.NoError(t, err)
	_, err = os.OpenFile(filepath.Join(testDir, "Scripts", "tool1"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	executable := env.FindExecutable("tool1")
	if assert.NotNil(t, executable) {
		assert.Equal(t, "tool1", executable.Path)
	}
}

func TestPython_Environment_FindExecutable_Failed(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	testDir := t.TempDir()

	env, err := tool.CreateEnvironment(".", testDir)
	assert.NoError(t, err)

	executable := env.FindExecutable("tool1")
	assert.Nil(t, executable)
}

func TestPython_Environment_RunExecutable_NotFound(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	testDir := t.TempDir()

	env, err := tool.CreateEnvironment(".", testDir)
	assert.NoError(t, err)

	err = env.RunExecutable(context.Background(), internal.RunOptions{}, "echo", []string{"arg1", "arg2"})
	assert.EqualError(t, err, "executable not found")
}
