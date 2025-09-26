package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func runWithEnvironment(ctx context.Context, environment internal.Environment, args ...string) error {
	return test.RunCommand(ctx, NewEnvironmentCommand(), internal.Context{
		Environment: environment,
	}, args...)
}

func TestEnvironment_List_Empty(t *testing.T) {
	env := test.NewEnvironment(t)

	err := runWithEnvironment(context.Background(), env, "list")

	require.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Status_Running(t *testing.T) {
	env := test.NewEnvironment(t)

	env.On("Id").Return("test")
	env.On("IsRunning", mock.Anything).Return(true)

	err := runWithEnvironment(context.Background(), env, "status")

	require.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Status_NotRunning(t *testing.T) {
	env := test.NewEnvironment(t)

	env.On("Id").Return("test")
	env.On("IsRunning", mock.Anything).Return(false)

	err := runWithEnvironment(context.Background(), env, "status")

	require.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Start(t *testing.T) {
	env := test.NewEnvironment(t)

	env.On("Start", mock.Anything).Return(nil)

	err := runWithEnvironment(context.Background(), env, "start")

	require.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Start_Failed(t *testing.T) {
	env := test.NewEnvironment(t)

	env.On("Start", mock.Anything).Return(errors.New("failed"))

	err := runWithEnvironment(context.Background(), env, "start")

	require.EqualError(t, err, "environment start failed: failed")
	env.AssertExpectations(t)
}

func TestEnvironment_Stop(t *testing.T) {
	env := test.NewEnvironment(t)

	env.On("Stop", mock.Anything).Return(nil)

	err := runWithEnvironment(context.Background(), env, "stop")

	require.NoError(t, err)
	env.AssertExpectations(t)
}

func TestEnvironment_Stop_Failed(t *testing.T) {
	env := test.NewEnvironment(t)

	env.On("Stop", mock.Anything).Return(errors.New("failed"))

	err := runWithEnvironment(context.Background(), env, "stop")

	require.EqualError(t, err, "environment stop failed: failed")
	env.AssertExpectations(t)
}
