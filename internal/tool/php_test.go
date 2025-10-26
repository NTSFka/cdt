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

func TestPHP_DetectPHP(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php").
		Return(env.NewExecutable("/bin/php"), nil)

	php := tool.DetectPHP(t.Context(), internal.ConfigTools{}, env)
	assert.NotNil(t, php)
	assert.Equal(t, "php", php.Id())
	assert.True(t, php.IsAvailable())

	if executable := php.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/php", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHP_DetectPHP_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php").
		Return(nil, nil)

	php := tool.DetectPHP(t.Context(), internal.ConfigTools{}, env)
	assert.NotNil(t, php)
	assert.Equal(t, "php", php.Id())
	assert.False(t, php.IsAvailable())
	assert.Nil(t, php.Executable())

	env.AssertExpectations(t)
}

func TestPHP_DetectPHP_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("php-8.2").
		Return(env.NewExecutable("/bin/php"), nil)

	php := tool.DetectPHP(t.Context(), internal.ConfigTools{
		"php": "php-8.2",
	}, env)
	assert.NotNil(t, php)
	assert.Equal(t, "php", php.Id())
	assert.True(t, php.IsAvailable())

	if executable := php.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/php", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHP_PHP_RunTarget(t *testing.T) {
	exec := test.NewExecutable(t)

	php := tool.NewPHP(exec.LazyExecutable("php"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("php", []string{"-f", "index.php"}).
		Return(nil)

	err := php.RunTarget(t.Context(), internal.ProjectRunnerOptions{ProjectInfo: info}, "index.php")
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHP_PHP_RunTarget_Fail(t *testing.T) {
	exec := test.NewExecutable(t)

	php := tool.NewPHP(exec.LazyExecutable("php"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("php", []string{"-f", "index.php"}).
		Return(errors.New("failed"))

	err := php.RunTarget(t.Context(), internal.ProjectRunnerOptions{ProjectInfo: info}, "index.php")
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
