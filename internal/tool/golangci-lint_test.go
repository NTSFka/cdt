package tool_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(env.NewExecutable("/bin/golangci-lint"))

	golangCILint := tool.DetectGolangCILint(t.Context(), env)
	assert.NotNil(t, golangCILint)
	assert.Equal(t, "golangci-lint", golangCILint.Id())
	assert.True(t, golangCILint.IsAvailable())

	if executable := golangCILint.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/golangci-lint", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestGolangCILint_DetectGolangCILint_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(nil)

	golangCILint := tool.DetectGolangCILint(t.Context(), env)
	assert.NotNil(t, golangCILint)
	assert.Equal(t, "golangci-lint", golangCILint.Id())
	assert.False(t, golangCILint.IsAvailable())
	assert.Nil(t, golangCILint.Executable())

	env.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run"}).
		Return(nil)

	err := golangCILint.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run", "mod1"}).
		Return(nil)

	err := golangCILint.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"mod1"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
