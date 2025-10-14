package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBandit_DetectBandit(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("bandit").
		Return(env.NewExecutable("/bin/bandit"), nil)

	bandit := tool.DetectBandit(t.Context(), env)
	assert.NotNil(t, bandit)
	assert.Equal(t, "bandit", bandit.Id())
	assert.True(t, bandit.IsAvailable())

	if executable := bandit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/bandit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestBandit_DetectBandit_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("bandit").
		Return(nil, nil)

	bandit := tool.DetectBandit(t.Context(), env)
	assert.NotNil(t, bandit)
	assert.Equal(t, "bandit", bandit.Id())
	assert.False(t, bandit.IsAvailable())
	assert.Nil(t, bandit.Executable())

	env.AssertExpectations(t)
}

func TestBandit_Bandit_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	bandit := tool.NewBandit(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"*"}).
		Return(nil)

	err := bandit.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBandit_Bandit_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	bandit := tool.NewBandit(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"file.py", "/path/to/file2.py"}).
		Return(nil)

	err := bandit.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"file.py", "/path/to/file2.py"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
