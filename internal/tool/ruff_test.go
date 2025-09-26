package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuff_DetectRuff(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("ruff").
		Return(env.NewExecutable("/bin/ruff"))

	ruff := tool.DetectRuff(t.Context(), env)
	assert.NotNil(t, ruff)
	assert.Equal(t, "ruff", ruff.Id())
	assert.True(t, ruff.IsAvailable())

	if executable := ruff.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/ruff", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestRuff_DetectRuff_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("ruff").
		Return(nil)

	ruff := tool.DetectRuff(t.Context(), env)
	assert.NotNil(t, ruff)
	assert.Equal(t, "ruff", ruff.Id())
	assert.False(t, ruff.IsAvailable())
	assert.Nil(t, ruff.Executable())

	env.AssertExpectations(t)
}

func TestRuff_Ruff_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"check"}).
		Return(nil)

	err := ruff.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"check", "file.py", "/path/to/file2.py"}).
		Return(nil)

	err := ruff.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"file.py", "/path/to/file2.py"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format"}).
		Return(nil)

	err := ruff.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "tests/*"}).
		Return(nil)

	err := ruff.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"tests/*"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "--check"}).
		Return(nil)

	err := ruff.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestRuff_Ruff_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	ruff := tool.NewRuff(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"format", "--check", "tests/*", "/path/to/file.py"}).
		Return(nil)

	err := ruff.FormatCheckFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"tests/*", "/path/to/file.py"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
