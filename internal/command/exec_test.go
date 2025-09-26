package command_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/command"
	"cdt/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func execRun(ctx context.Context, environment internal.Environment, args ...string) error {
	return test.RunCommand(ctx, command.NewExecCommand(), internal.Context{
		Environment: environment,
	}, args...)
}

func TestExec_NoCommand(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	err := execRun(t.Context(), env)

	require.EqualError(t, err, "COMMAND is required")

	env.AssertExpectations(t)
}

func TestExec_Target_Success(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	env.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"Hello!"}).
		Return(nil)

	err := execRun(t.Context(), env, "echo", "Hello!")
	require.NoError(t, err)

	env.AssertExpectations(t)
}

func TestExec_Target_Failure(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	env.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"Hello!"}).
		Return(errors.New("failed"))

	err := execRun(t.Context(), env, "echo", "Hello!")

	require.Error(t, err)
	assert.Equal(t, "command failed: failed", err.Error())

	env.AssertExpectations(t)
}
