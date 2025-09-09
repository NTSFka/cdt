package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPip_DetectPip_Pip(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pip").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPip(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pip", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPip_DetectPip_System(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pip").
		Return(nil)

	env.OnFindExecutable("pip3").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPip(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pip", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPip_DetectPip_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pip").
		Return(nil)
	env.OnFindExecutable("pip3").
		Return(nil)

	tool := DetectPip(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pip", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPip_Pip_AddDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPip(exec.LazyExecutable("pip"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("pip", []string{"install", "dep1"}).
		Return(nil)

	err := tool.AddDependencies(desc, []string{"dep1"}, false)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_RemoveDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPip(exec.LazyExecutable("pip"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("pip", []string{"uninstall", "dep1"}).
		Return(nil)

	err := tool.RemoveDependencies(desc, []string{"dep1"}, false)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_UpdateDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPip(exec.LazyExecutable("pip"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("pip", []string{"install", "--upgrade", "dep1"}).
		Return(nil)

	err := tool.UpdateDependencies(desc, []string{"dep1"})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_FetchDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPip(exec.LazyExecutable("pip"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("pip", []string{"install", "-r", "requirements.txt"}).
		Return(nil)

	err := tool.FetchDependencies(desc, false)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_ListDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPip(exec.LazyExecutable("pip"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("pip", []string{"list"}).
		Return(nil)

	err := tool.ListDependencies(desc)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_AuditDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPip(exec.LazyExecutable("pip"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("pip", []string{"audit"}).
		Return(nil)

	err := tool.AuditDependencies(desc)
	assert.EqualError(t, err, "not supported")

	exec.AssertExpectations(t)
}
