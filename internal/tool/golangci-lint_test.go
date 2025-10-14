package tool_test

import (
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGolangCILint_DetectGolangCILint(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("golangci-lint").
		Return(env.NewExecutable("/bin/golangci-lint"), nil)

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
		Return(nil, nil)

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

	err := golangCILint.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"mod1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt"}).
		Return(nil)

	err := golangCILint.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt"}).
		Return(errors.New("failed"))

	err := golangCILint.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "file1"}).
		Return(nil)

	err := golangCILint.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"file1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "file1"}).
		Return(errors.New("failed"))

	err := goTool.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"file1"},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff"}).
		Return(nil)

	err := golangCILint.FormatCheckAll(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatCheckAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff"}).
		Return(errors.New("failed"))

	err := golangCILint.FormatCheckAll(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	golangCILint := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff", "file1"}).
		Return(nil)

	err := golangCILint.FormatCheckFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"file1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGolangCILint_FormatCheckFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGolangCILint(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"fmt", "--diff", "file1"}).
		Return(errors.New("failed"))

	err := goTool.FormatCheckFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"file1"},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
