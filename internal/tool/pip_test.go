package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPip_DetectPip_Pip(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pip").
		Return(env.NewExecutable("/bin/pip"), nil)

	pip := tool.DetectPip(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, pip)
	assert.Equal(t, "pip", pip.Id())
	assert.True(t, pip.IsAvailable())

	if executable := pip.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pip", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPip_DetectPip_System(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pip").
		Return(nil, nil)

	env.OnFindExecutable("pip3").
		Return(env.NewExecutable("/bin/pip"), nil)

	pip := tool.DetectPip(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, pip)
	assert.Equal(t, "pip", pip.Id())
	assert.True(t, pip.IsAvailable())

	if executable := pip.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pip", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPip_DetectPip_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pip").
		Return(nil, nil)
	env.OnFindExecutable("pip3").
		Return(nil, nil)

	pip := tool.DetectPip(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, pip)
	assert.Equal(t, "pip", pip.Id())
	assert.False(t, pip.IsAvailable())
	assert.Nil(t, pip.Executable())

	env.AssertExpectations(t)
}

func TestPip_DetectPip_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pip3").
		Return(env.NewExecutable("/bin/pip"), nil)

	pip := tool.DetectPip(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"pip": "pip3"},
	})
	assert.NotNil(t, pip)
	assert.Equal(t, "pip", pip.Id())
	assert.True(t, pip.IsAvailable())

	if executable := pip.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pip", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPip_Pip_AddDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	pip := tool.NewPip(exec.LazyExecutable("pip"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("pip", []string{"install", "dep1"}).
		Return(nil)

	err := pip.AddDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_RemoveDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	pip := tool.NewPip(exec.LazyExecutable("pip"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("pip", []string{"uninstall", "dep1"}).
		Return(nil)

	err := pip.RemoveDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_UpdateDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	pip := tool.NewPip(exec.LazyExecutable("pip"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("pip", []string{"install", "--upgrade", "dep1"}).
		Return(nil)

	err := pip.UpdateDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_FetchDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	pip := tool.NewPip(exec.LazyExecutable("pip"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("pip", []string{"install", "-r", "requirements.txt"}).
		Return(nil)

	err := pip.FetchDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_ListDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	pip := tool.NewPip(exec.LazyExecutable("pip"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("pip", []string{"list"}).
		Return(nil)

	err := pip.ListDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPip_Pip_AuditDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	pip := tool.NewPip(exec.LazyExecutable("pip"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("pip", []string{"audit"}).
		Return(nil)

	err := pip.AuditDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
	)
	require.EqualError(t, err, "not supported")

	exec.AssertExpectations(t)
}
