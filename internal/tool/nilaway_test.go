package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilAway_DetectNilAway(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("nilaway").
		Return(env.NewExecutable("/bin/nilaway"), nil)

	nilAway := tool.DetectNilAway(t.Context(), tool.DetectOptions{Environment: env})
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
		Return(nil, nil)

	nilAway := tool.DetectNilAway(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, nilAway)
	assert.Equal(t, "nilaway", nilAway.Id())
	assert.False(t, nilAway.IsAvailable())
	assert.Nil(t, nilAway.Executable())

	env.AssertExpectations(t)
}

func TestNilAway_DetectNilAway_Config(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("nilaway-2").
		Return(env.NewExecutable("/bin/nilaway"), nil)

	nilAway := tool.DetectNilAway(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"nilaway": "nilaway-2"},
	})
	assert.NotNil(t, nilAway)
	assert.Equal(t, "nilaway", nilAway.Id())
	assert.True(t, nilAway.IsAvailable())

	if executable := nilAway.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/nilaway", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestNilAway_NilAway_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	nilAway := tool.NewNilAway(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"./..."}).
		Return(nil)

	err := nilAway.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestNilAway_NilAway_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	nilAway := tool.NewNilAway(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"mod1"}).
		Return(nil)

	err := nilAway.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info, Filenames: &[]string{"mod1"}},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
