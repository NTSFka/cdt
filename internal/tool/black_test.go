package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlack_DetectBlack(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(env.NewExecutable("/bin/black"))

	black := tool.DetectBlack(t.Context(), env)
	assert.NotNil(t, black)
	assert.Equal(t, "black", black.Id())
	assert.True(t, black.IsAvailable())

	if executable := black.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/black", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestBlack_DetectBlack_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(nil)

	black := tool.DetectBlack(t.Context(), env)
	assert.NotNil(t, black)
	assert.Equal(t, "black", black.Id())
	assert.False(t, black.IsAvailable())
	assert.Nil(t, black.Executable())

	env.AssertExpectations(t)
}

func TestBlack_Black_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{}).
		Return(nil)

	err := black.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"tests/*"}).
		Return(nil)

	err := black.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"tests/*"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"--check"}).
		Return(nil)

	err := black.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	black := tool.NewBlack(exec.LazyExecutable("format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"--check", "tests/*", "/path/to/file.py"}).
		Return(nil)

	err := black.FormatCheckFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"tests/*", "/path/to/file.py"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
