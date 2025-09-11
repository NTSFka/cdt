package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuff_DetectRuff(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("ruff").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectRuff(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "ruff", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestRuff_DetectRuff_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("ruff").
		Return(nil)

	tool := DetectRuff(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "ruff", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestRuff_Ruff_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewRuff(exec.LazyExecutable("lint"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"check"}).
		Return(nil)

	err := tool.LintAll(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewRuff(exec.LazyExecutable("lint"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"check", "file.py", "/path/to/file2.py"}).
		Return(nil)

	err := tool.LintFiles(context.Background(), desc, []string{"file.py", "/path/to/file2.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format"}).
		Return(nil)

	err := tool.FormatAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "tests/*"}).
		Return(nil)

	err := tool.FormatFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"tests/*"})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "--check"}).
		Return(nil)

	err := tool.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "--check", "tests/*", "/path/to/file.py"}).
		Return(nil)

	err := tool.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"tests/*", "/path/to/file.py"})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
