package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*"}).
		Return(nil)

	err := tool.LintAll(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBandit_Bandit_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBandit(exec.LazyExecutable("lint"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(context.Background(), desc, []string{"file.py", "/path/to/file2.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
