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

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
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
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestPython_Python_RunTarget(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPython(exec.LazyExecutable("python3"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Runner: tool})

	exec.OnRun("python3", []string{"main.py"}).
		Return(nil)

	err := tool.RunTarget(p, "main.py", []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPython_Python_RunTarget_Fail(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPython(exec.LazyExecutable("python3"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Runner: tool})

	exec.OnRun("python3", []string{"main.py"}).
		Return(errors.New("failed"))

	err := tool.RunTarget(p, "main.py", []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestPython_CreateEnvironment(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)
	assert.Equal(t, "pyenv", tool.IdShort())

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)
	assert.NotNil(t, env)

	assert.Equal(t, "pyenv", env.Id())
}

func TestPython_Environment_Start(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	err = env.Start(context.Background())
	assert.NoError(t, err)
}

func TestPython_Environment_IsRunning_True(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	result := env.IsRunning(context.Background())
	assert.True(t, result)
}

func TestPython_Environment_Stop(t *testing.T) {

	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	err = env.Stop(context.Background())
	assert.NoError(t, err)
}

func TestPython_Environment_Cleanup(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	env, err := tool.CreateEnvironment(".", ".venv")
	assert.NoError(t, err)

	err = env.Cleanup(context.Background())
	assert.NoError(t, err)
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
		assert.Equal(t, filepath.Join(testDir, "bin", "tool1"), executable.Path)
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
		assert.Equal(t, filepath.Join(testDir, "Scripts", "tool1"), executable.Path)
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

func TestPython_Environment_RunExecutable(t *testing.T) {
	tool := NewPython(func() *internal.Executable { return &internal.Executable{Path: "python"} })
	assert.NotNil(t, tool)

	testDir := t.TempDir()

	env, err := tool.CreateEnvironment(".", testDir)
	assert.NoError(t, err)

	err = env.RunExecutable(context.Background(), internal.RunOptions{}, "echo", []string{"arg1", "arg2"})
	assert.NoError(t, err)
}
