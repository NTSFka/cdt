package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMyPy_DetectMyPy(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(env.NewExecutable("/bin/mypy"), nil)

	mypy := tool.DetectMyPy(t.Context(), env)
	assert.NotNil(t, mypy)
	assert.Equal(t, "mypy", mypy.Id())
	assert.True(t, mypy.IsAvailable())

	if executable := mypy.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/mypy", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestMyPy_DetectMyPy_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("mypy").
		Return(nil, nil)

	mypy := tool.DetectMyPy(t.Context(), env)
	assert.NotNil(t, mypy)
	assert.Equal(t, "mypy", mypy.Id())
	assert.False(t, mypy.IsAvailable())
	assert.Nil(t, mypy.Executable())

	env.AssertExpectations(t)
}

func TestMyPy_MyPy_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	mypy := tool.NewMyPy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*.py"}).
		Return(nil)

	err := mypy.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestMyPy_MyPy_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	mypy := tool.NewMyPy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := mypy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"file.py", "/path/to/file2.py"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
