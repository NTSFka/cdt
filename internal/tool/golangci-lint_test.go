package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetectGolangCILint(t *testing.T) {
	environment := test.Environment{}
	environment.On("FindExecutable", "golangci-lint").Return(environment.MakeExecutable("/bin/tool"))

	tool := DetectGolangCILint(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "golangci-lint", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}

	environment.AssertExpectations(t)
}

func TestDetectGolangCILintNotFound(t *testing.T) {
	environment := test.Environment{}
	environment.On("FindExecutable", "golangci-lint").Return(nil)

	tool := DetectGolangCILint(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "golangci-lint", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestGolangCILintLintAll(t *testing.T) {
	environment := test.Environment{}

	tool := NewGolangCILint(environment.MakeExecutable("lint"))

	p := internal.MakeProject("project", "", nil, internal.Workflow{Linter: tool})

	environment.OnRunSuccess(p, "lint", []string{"run"})

	err := tool.LintAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestGolangCILintLint(t *testing.T) {
	environment := test.Environment{}

	tool := NewGolangCILint(environment.MakeExecutable("lint"))

	p := internal.MakeProject("project", "", nil, internal.Workflow{Linter: tool})

	environment.OnRunSuccess(p, "lint", []string{"run", "mod1"})

	err := tool.LintFiles(p, []string{"mod1"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}
