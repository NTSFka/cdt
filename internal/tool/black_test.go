package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBlack_DetectBlack(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectBlack(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "black", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}

	env.AssertExpectations(t)
}

func TestBlack_DetectBlack_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("black").
		Return(nil)

	tool := DetectBlack(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "black", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	env.AssertExpectations(t)
}

func TestBlack_Black_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Formatter: tool})

	exec.OnRun("format", []string{}).
		Return(nil)

	err := tool.FormatAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Formatter: tool})

	exec.OnRun("format", []string{"tests/*"}).
		Return(nil)

	err := tool.FormatFiles(p, []string{"tests/*"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Formatter: tool})

	exec.OnRun("format", []string{"--check"}).
		Return(nil)

	err := tool.FormatCheckAll(p, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestBlack_Black_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewBlack(exec.LazyExecutable("format"))

	p := internal.MakeProject(".", "", nil, internal.Workflow{Formatter: tool})

	exec.OnRun("format", []string{"--check", "tests/*", "/path/to/file.py"}).
		Return(nil)

	err := tool.FormatCheckFiles(p, []string{"tests/*", "/path/to/file.py"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}
