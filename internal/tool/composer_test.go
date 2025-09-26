package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComposer_DetectComposer_Phar(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("composer.phar").
		Return(env.NewExecutable("/bin/composer"))

	composer := tool.DetectComposer(t.Context(), env)
	assert.NotNil(t, composer)
	assert.Equal(t, "composer", composer.Id())
	assert.True(t, composer.IsAvailable())

	if executable := composer.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/composer", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestComposer_DetectComposer_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("composer.phar").
		Return(nil)

	// System installation
	env.OnFindExecutable("composer").
		Return(env.NewExecutable("/bin/composer"))

	composer := tool.DetectComposer(t.Context(), env)
	assert.NotNil(t, composer)
	assert.Equal(t, "composer", composer.Id())
	assert.True(t, composer.IsAvailable())

	if executable := composer.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/composer", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestComposer_DetectComposer_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("composer.phar").
		Return(nil)
	env.OnFindExecutable("composer").
		Return(nil)

	composer := tool.DetectComposer(t.Context(), env)
	assert.NotNil(t, composer)
	assert.Equal(t, "composer", composer.Id())
	assert.False(t, composer.IsAvailable())
	assert.Nil(t, composer.Executable())

	env.AssertExpectations(t)
}

func TestComposer_Composer_AddDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"require", "dep1"}).
		Return(nil)

	err := composer.AddDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_AddDependencies_Dev(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"require", "--dev", "dep1"}).
		Return(nil)

	err := composer.AddDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		true,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_RemoveDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"remove", "dep1"}).
		Return(nil)

	err := composer.RemoveDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_RemoveDependencies_Dev(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"remove", "--dev", "dep1"}).
		Return(nil)

	err := composer.RemoveDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		true,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_UpdateDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"update", "dep1"}).
		Return(nil)

	err := composer.UpdateDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_FetchDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"install"}).
		Return(nil)

	err := composer.FetchDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_FetchDependencies_NoDev(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"install", "--no-dev"}).
		Return(nil)

	err := composer.FetchDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		true,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_ListDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"show"}).
		Return(nil)

	err := composer.ListDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_AuditDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	composer := tool.NewComposer(exec.LazyExecutable("composer"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("composer", []string{"audit"}).
		Return(nil)

	err := composer.AuditDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
