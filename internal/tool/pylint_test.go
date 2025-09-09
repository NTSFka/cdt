package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPylint_DetectPylint(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pylint").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPylint(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pylint", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPylint_DetectPylint_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pylint").
		Return(nil)

	tool := DetectPylint(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "pylint", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPylint_Pylint_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPylint(exec.LazyExecutable("lint"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("lint", []string{"*"}).
		Return(nil)

	err := tool.LintAll(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPylint_Pylint_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPylint(exec.LazyExecutable("lint"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(desc, []string{"file.py", "/path/to/file2.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
