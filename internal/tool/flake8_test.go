package tool_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlake8_DetectFlake8(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8").
		Return(env.NewExecutable("/bin/flake8"))

	flake8 := tool.DetectFlake8(t.Context(), env)
	assert.NotNil(t, flake8)
	assert.Equal(t, "flake8", flake8.Id())
	assert.True(t, flake8.IsAvailable())

	if executable := flake8.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/flake8", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestFlake8_DetectFlake8_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8").
		Return(nil)

	flake8 := tool.DetectFlake8(t.Context(), env)
	assert.NotNil(t, flake8)
	assert.Equal(t, "flake8", flake8.Id())
	assert.False(t, flake8.IsAvailable())
	assert.Nil(t, flake8.Executable())

	env.AssertExpectations(t)
}

func TestFlake8_Flake8_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	flake8 := tool.NewFlake8(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{}).
		Return(nil)

	err := flake8.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestFlake8_Flake8_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	flake8 := tool.NewFlake8(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := flake8.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"file.py", "/path/to/file2.py"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
