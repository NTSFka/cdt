package test

import (
	"context"
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/mock"
)

func LazyExecutableNil() (*internal.Executable, error) {
	return nil, nil // nolint: nilnil
}

func LazyExecutable(path string) func() (*internal.Executable, error) {
	return func() (*internal.Executable, error) {
		return &internal.Executable{Path: path}, nil
	}
}

// Executable allow testing executable invocation.
type Executable struct {
	mock.Mock
	Runtime ExecutableRuntime
}

// NewExecutable create new testing executable.
func NewExecutable(t *testing.T) *Executable {
	executable := Executable{}
	executable.Test(t)
	executable.Runtime.Test(t)
	executable.Runtime.On("Id").Return("test")

	return &executable
}

// OnRun set an expectation on calling executable.
func (m *Executable) OnRun(path string, args []string) *mock.Call {
	return m.Runtime.On("RunExecutable", mock.Anything, mock.Anything, path, args)
}

// OnRunAnything set an expectation on calling executable without any arguments.
func (m *Executable) OnRunAnything(path string) *mock.Call {
	return m.Runtime.On("RunExecutable", mock.Anything, mock.Anything, path, mock.Anything)
}

// OnRunOutput set an expectation on calling executable and printing output to stdout.
func (m *Executable) OnRunOutput(path string, args []string, output string) *mock.Call {
	return m.OnRun(path, args).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(internal.RunOptions)
			_, _ = c.Output.Write([]byte(output))
		})
}

// NewExecutable creates a new executable.
func (m *Executable) NewExecutable(path string) *internal.Executable {
	return &internal.Executable{Path: path, Runtime: &m.Runtime}
}

// LazyExecutable creates a new lazy executable via function call.
func (m *Executable) LazyExecutable(path string) func() (*internal.Executable, error) {
	return func() (*internal.Executable, error) {
		return m.NewExecutable(path), nil
	}
}

type ExecutableRuntime struct {
	mock.Mock
}

func NewExecutableRuntime(t *testing.T) *ExecutableRuntime {
	runtime := ExecutableRuntime{}
	runtime.Test(t)

	return &runtime
}

func (t *ExecutableRuntime) Id() string {
	return t.Called().Get(0).(string)
}

func (t *ExecutableRuntime) RunExecutable(
	ctx context.Context,
	options internal.RunOptions,
	path string,
	args []string,
) error {
	return t.Called(ctx, options, path, args).Error(0)
}
