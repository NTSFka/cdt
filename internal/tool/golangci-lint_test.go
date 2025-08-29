package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectGolangCILint(env)
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

	tool := DetectGolangCILint(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "golangci-lint", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGolangCILint(exec.LazyExecutable("lint"))

	p := internal.MakeProject("project", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"run"}).
		Return(nil)

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_GolangCILint_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewGolangCILint(exec.LazyExecutable("lint"))

	p := internal.MakeProject("project", "", nil, internal.Workflow{Linter: tool})

	exec.OnRun("lint", []string{"run", "mod1"}).
		Return(nil)

	err := tool.LintFiles(p, []string{"mod1"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
