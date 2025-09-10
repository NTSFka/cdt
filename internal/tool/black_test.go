package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlack_DetectBlack(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectBlack(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "black", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestBlack_DetectBlack_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(nil)

	tool := DetectBlack(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "black", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestBlack_Black_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{}).
		Return(nil)

	err := tool.FormatAll(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"tests/*"}).
		Return(nil)

	err := tool.FormatFiles(context.Background(), desc, []string{"tests/*"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"--check"}).
		Return(nil)

	err := tool.FormatCheckAll(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("format", []string{"--check", "tests/*", "/path/to/file.py"}).
		Return(nil)

	err := tool.FormatCheckFiles(context.Background(), desc, []string{"tests/*", "/path/to/file.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
