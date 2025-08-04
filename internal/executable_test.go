package internal

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type testExecutableRuntime struct {
	mock.Mock
}

func (t *testExecutableRuntime) Id() string {
	return t.Called().Get(0).(string)
}

func (t *testExecutableRuntime) RunExecutable(ctx context.Context, options RunOptions, path string, args []string) error {
	return t.Called(ctx, options, path, args).Error(0)
}

func TestExecutable_Run(t *testing.T) {
	runtime := &testExecutableRuntime{}

	executable := Executable{Path: "echo", Runtime: runtime}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{}).
		Return(nil)

	err := executable.Run(context.Background(), RunOptions{}, []string{})
	assert.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestExecutable_Run_Failed(t *testing.T) {
	runtime := &testExecutableRuntime{}

	executable := Executable{Path: "echo", Runtime: runtime}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{}).
		Return(errors.New("failed"))

	err := executable.Run(context.Background(), RunOptions{}, []string{})
	assert.EqualError(t, err, "failed")

	runtime.AssertExpectations(t)
}

func TestExecutable_Run_Args(t *testing.T) {
	runtime := &testExecutableRuntime{}

	executable := Executable{Path: "echo", Runtime: runtime}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "echo", []string{"arg1", "arg2"}).
		Return(nil)

	err := executable.Run(context.Background(), RunOptions{}, []string{"arg1", "arg2"})
	assert.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestExecutable_Run_ArgsExtra(t *testing.T) {
	runtime := &testExecutableRuntime{}

	executable := Executable{
		Path:    "print",
		Args:    []string{"arg1", "arg2"},
		Runtime: runtime,
	}

	runtime.On("RunExecutable", mock.Anything, mock.Anything, "print", []string{"arg1", "arg2", "arg3", "arg4"}).
		Return(nil)

	err := executable.Run(context.Background(), RunOptions{}, []string{"arg3", "arg4"})
	assert.NoError(t, err)

	runtime.AssertExpectations(t)
}
