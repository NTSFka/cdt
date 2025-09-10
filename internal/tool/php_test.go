package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPHP_DetectPHP(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php").
		Return(env.NewExecutable("/bin/tool"))

	tool := DetectPHP(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "php", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHP_DetectPHP_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php").
		Return(nil)

	tool := DetectPHP(env)
	assert.NotNil(t, tool)
	assert.Equal(t, "php", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestPHP_PHP_RunTarget(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHP(exec.LazyExecutable("php"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("php", []string{"-f", "index.php"}).
		Return(nil)

	err := tool.RunTarget(context.Background(), desc, "index.php", []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHP_PHP_RunTarget_Fail(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewPHP(exec.LazyExecutable("php"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("php", []string{"-f", "index.php"}).
		Return(errors.New("failed"))

	err := tool.RunTarget(context.Background(), desc, "index.php", []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
