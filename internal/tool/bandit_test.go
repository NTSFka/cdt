package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBandit_DetectBandit(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("bandit").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectBandit(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "bandit", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestBandit_DetectBandit_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("bandit").
		Return(nil)

	tool := DetectBandit(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "bandit", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestBandit_Bandit_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBandit(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*"}).
		Return(nil)

	err := tool.LintAll(context.Background(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBandit_Bandit_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBandit(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(context.Background(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"file.py", "/path/to/file2.py"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
