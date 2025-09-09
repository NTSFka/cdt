package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlake8_DetectFlake8(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectFlake8(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "flake8", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestFlake8_DetectFlake8_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8").
		Return(nil)

	tool := DetectFlake8(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "flake8", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestFlake8_Flake8_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewFlake8(exec.LazyExecutable("lint"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("lint", []string{}).
		Return(nil)

	err := tool.LintAll(desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestFlake8_Flake8_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewFlake8(exec.LazyExecutable("lint"))

	desc := internal.Project{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(desc, []string{"file.py", "/path/to/file2.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
