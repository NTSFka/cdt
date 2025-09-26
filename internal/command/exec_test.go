package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func execRun(ctx context.Context, environment internal.Environment, args ...string) error {
	return test.RunCommand(ctx, NewExecCommand(), internal.Context{
		Environment: environment,
	}, args...)
}

func TestExec_NoCommand(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	err := execRun(context.Background(), env)

	require.EqualError(t, err, "COMMAND is required")

	env.AssertExpectations(t)
}

func TestExec_Target_Success(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	env.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"Hello!"}).
		Return(nil)

	err := execRun(context.Background(), env, "echo", "Hello!")
	require.NoError(t, err)

	env.AssertExpectations(t)
}

func TestExec_Target_Failure(t *testing.T) {
	env := test.NewEnvironment(t)
	env.Test(t)

	env.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"Hello!"}).
		Return(errors.New("failed"))

	err := execRun(context.Background(), env, "echo", "Hello!")

	require.Error(t, err)
	assert.Equal(t, "command failed: failed", err.Error())

	env.AssertExpectations(t)
}
