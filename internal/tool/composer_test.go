package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposer_DetectComposer_Phar(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("composer.phar").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectComposer(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "composer", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
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
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectComposer(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "composer", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestComposer_DetectComposer_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("composer.phar").
		Return(nil)
	env.OnFindExecutable("composer").
		Return(nil)

	tool := DetectComposer(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "composer", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestComposer_Composer_AddDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"require", "dep1"}).
		Return(nil)

	err := tool.AddDependencies(desc, []string{"dep1"}, false)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_AddDependencies_Dev(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"require", "--dev", "dep1"}).
		Return(nil)

	err := tool.AddDependencies(desc, []string{"dep1"}, true)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_RemoveDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"remove", "dep1"}).
		Return(nil)

	err := tool.RemoveDependencies(desc, []string{"dep1"}, false)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_RemoveDependencies_Dev(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"remove", "--dev", "dep1"}).
		Return(nil)

	err := tool.RemoveDependencies(desc, []string{"dep1"}, true)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_UpdateDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"update", "dep1"}).
		Return(nil)

	err := tool.UpdateDependencies(desc, []string{"dep1"})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_FetchDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"install"}).
		Return(nil)

	err := tool.FetchDependencies(desc, false)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_FetchDependencies_NoDev(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"install", "--no-dev"}).
		Return(nil)

	err := tool.FetchDependencies(desc, true)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_ListDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"show"}).
		Return(nil)

	err := tool.ListDependencies(desc)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestComposer_Composer_AuditDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewComposer(exec.LazyExecutable("composer"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("composer", []string{"audit"}).
		Return(nil)

	err := tool.AuditDependencies(desc)
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
