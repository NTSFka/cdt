package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectGolangCILint(context.Background(), env)
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

	tool := DetectGolangCILint(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "golangci-lint", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGolangCILint(exec.LazyExecutable("lint"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run"}).
		Return(nil)

	err := tool.LintAll(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGolangCILint(exec.LazyExecutable("lint"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"run", "mod1"}).
		Return(nil)

	err := tool.LintFiles(context.Background(), desc, []string{"mod1"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
