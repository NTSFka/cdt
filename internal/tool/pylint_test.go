package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPylint_DetectPylint(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pylint").
		Return(env.NewExecutable("/bin/pylint"), nil)

	pylint := tool.DetectPylint(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, pylint)
	assert.Equal(t, "pylint", pylint.Id())
	assert.True(t, pylint.IsAvailable())

	if executable := pylint.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pylint", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPylint_DetectPylint_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pylint").
		Return(nil, nil)

	pylint := tool.DetectPylint(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, pylint)
	assert.Equal(t, "pylint", pylint.Id())
	assert.False(t, pylint.IsAvailable())
	assert.Nil(t, pylint.Executable())

	env.AssertExpectations(t)
}

func TestPylint_DetectPylint_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pylint4").
		Return(env.NewExecutable("/bin/pylint"), nil)

	pylint := tool.DetectPylint(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"pylint": "pylint4"},
	})
	assert.NotNil(t, pylint)
	assert.Equal(t, "pylint", pylint.Id())
	assert.True(t, pylint.IsAvailable())

	if executable := pylint.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pylint", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPylint_Pylint_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	pylint := tool.NewPylint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*"}).
		Return(nil)

	err := pylint.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPylint_Pylint_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	pylint := tool.NewPylint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := pylint.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file.py", "/path/to/file2.py"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
