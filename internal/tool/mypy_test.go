package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyPy_DetectMyPy(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectMyPy(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "mypy", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}

	env.AssertExpectations(t)
}

func TestMyPy_DetectMyPy_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(nil)

	tool := DetectMyPy(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "mypy", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestMyPy_MyPy_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewMyPy(exec.LazyExecutable("lint"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"*.py"}).
		Return(nil)

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestMyPy_MyPy_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewMyPy(exec.LazyExecutable("lint"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(p, []string{"file.py", "/path/to/file2.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
