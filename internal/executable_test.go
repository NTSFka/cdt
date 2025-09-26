package internal_test

import (
	"cdt/internal"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testExecutableRuntime struct {
	mock.Mock
}

func (t *testExecutableRuntime) Id() string {
	return t.Called().Get(0).(string)
}

func (t *testExecutableRuntime) RunExecutable(ctx context.Context, options internal.RunOptions, path string, args []string) error {
	return t.Called(ctx, options, path, args).Error(0)
}

func TestExecutable_Run(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.On("Id").Return("test")

	executable := internal.Executable{Path: "echo", Runtime: runtime}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{}).
		Return(nil)

	err := executable.Run(t.Context(), internal.RunOptions{}, []string{})
	require.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestExecutable_Run_Failed(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)
	runtime.On("Id").Return("test")

	executable := internal.Executable{Path: "echo", Runtime: runtime}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{}).
		Return(errors.New("failed"))

	err := executable.Run(t.Context(), internal.RunOptions{}, []string{})
	require.EqualError(t, err, "failed")

	runtime.AssertExpectations(t)
}

func TestExecutable_Run_Args(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)
	runtime.On("Id").Return("test")

	executable := internal.Executable{Path: "echo", Runtime: runtime}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"arg1", "arg2"}).
		Return(nil)

	err := executable.Run(t.Context(), internal.RunOptions{}, []string{"arg1", "arg2"})
	require.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestExecutable_Run_ArgsExtra(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)
	runtime.On("Id").Return("test")

	executable := internal.Executable{
		Path:    "print",
		Args:    []string{"arg1", "arg2"},
		Runtime: runtime,
	}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "print", []string{"arg1", "arg2", "arg3", "arg4"}).
		Return(nil)

	err := executable.Run(t.Context(), internal.RunOptions{}, []string{"arg3", "arg4"})
	require.NoError(t, err)

	runtime.AssertExpectations(t)
}
