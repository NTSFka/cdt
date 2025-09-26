package tool_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilAway_DetectNilAway(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("nilaway").
		Return(env.NewExecutable("/bin/nilaway"))

	nilAway := tool.DetectNilAway(t.Context(), env)
	assert.NotNil(t, nilAway)
	assert.Equal(t, "nilaway", nilAway.Id())
	assert.True(t, nilAway.IsAvailable())

	if executable := nilAway.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/nilaway", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestNilAway_DetectNilAway_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("nilaway").
		Return(nil)

	nilAway := tool.DetectNilAway(t.Context(), env)
	assert.NotNil(t, nilAway)
	assert.Equal(t, "nilaway", nilAway.Id())
	assert.False(t, nilAway.IsAvailable())
	assert.Nil(t, nilAway.Executable())

	env.AssertExpectations(t)
}

func TestNilAway_NilAway_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	nilAway := tool.NewNilAway(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"./..."}).
		Return(nil)

	err := nilAway.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestNilAway_NilAway_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	nilAway := tool.NewNilAway(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"mod1"}).
		Return(nil)

	err := nilAway.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"mod1"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
