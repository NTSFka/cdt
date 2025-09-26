package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlake8_DetectFlake8(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("flake8").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectFlake8(context.Background(), env)
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

	tool := DetectFlake8(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "flake8", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestFlake8_Flake8_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewFlake8(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{}).
		Return(nil)

	err := tool.LintAll(context.Background(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestFlake8_Flake8_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewFlake8(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(context.Background(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"file.py", "/path/to/file2.py"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
