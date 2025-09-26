package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyPy_DetectMyPy(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectMyPy(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "mypy", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestMyPy_DetectMyPy_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(nil)

	tool := DetectMyPy(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "mypy", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestMyPy_MyPy_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewMyPy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*.py"}).
		Return(nil)

	err := tool.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestMyPy_MyPy_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewMyPy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"file.py", "/path/to/file2.py"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
