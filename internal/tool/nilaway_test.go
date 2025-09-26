package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilAway_DetectNilAway(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("nilaway").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectNilAway(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "nilaway", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestNilAway_DetectNilAway_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("nilaway").
		Return(nil)

	tool := DetectNilAway(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "nilaway", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestNilAway_NilAway_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewNilAway(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"./..."}).
		Return(nil)

	err := tool.LintAll(context.Background(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestNilAway_NilAway_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewNilAway(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"mod1"}).
		Return(nil)

	err := tool.LintFiles(context.Background(), internal.ProjectLinterOptions{ProjectInfo: info}, []string{"mod1"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
