package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func execRun(environment internal.Environment, args ...string) error {
	return test.RunCommand(NewExecCommand(), internal.Context{
		Environment: environment,
	}, args...)
}

func TestExec_NoCommand(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	err := execRun(env)

	assert.EqualError(t, err, "COMMAND is required")

	env.AssertExpectations(t)
}

func TestExec_Target_Success(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	env.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"Hello!"}).
		Return(nil)

	err := execRun(env, "echo", "Hello!")

	assert.NoError(t, err)

	env.AssertExpectations(t)
}

func TestExec_Target_Failure(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	env.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"Hello!"}).
		Return(errors.New("failed"))

	err := execRun(env, "echo", "Hello!")

	if assert.Error(t, err) {
		assert.Equal(t, "command failed: failed", err.Error())
	}

	env.AssertExpectations(t)
}
