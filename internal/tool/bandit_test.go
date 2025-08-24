package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBandit_DetectBandit(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("bandit").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectBandit(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "bandit", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}

	env.AssertExpectations(t)
}

func TestBandit_DetectBandit_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("bandit").
		Return(nil)

	tool := DetectBandit(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "bandit", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestBandit_Bandit_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBandit(exec.LazyExecutable("lint"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"*"}).
		Return(nil)

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBandit_Bandit_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBandit(exec.LazyExecutable("lint"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(p, []string{"file.py", "/path/to/file2.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
