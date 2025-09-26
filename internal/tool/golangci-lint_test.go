package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectGolangCILint(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "golangci-lint", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestGolangCILint_DetectGolangCILint_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(nil)

	tool := DetectGolangCILint(t.Context(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "golangci-lint", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run"}).
		Return(nil)

	err := tool.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run", "mod1"}).
		Return(nil)

	err := tool.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"mod1"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
