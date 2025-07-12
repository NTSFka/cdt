package test

import (
	"cdt/internal"
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
)

// Executable allow testing executable invocation
type Executable struct {
	mock.Mock
}

// NewExecutable create new testing executable
func NewExecutable(t *testing.T) *Executable {
	executable := Executable{}
	executable.Test(t)
	return &executable
}

func (m *Executable) runFunc(ctx context.Context, options internal.RunOptions, path string, args []string) error {
	return m.Called(ctx, options, path, args).Error(0)
}

// OnRun set an expectation on calling executable
func (m *Executable) OnRun(path string, args []string) *mock.Call {
	return m.On("runFunc", mock.Anything, mock.Anything, path, args)
}

// OnRunAnything set an expectation on calling executable without any arguments
func (m *Executable) OnRunAnything(path string) *mock.Call {
	return m.On("runFunc", mock.Anything, mock.Anything, path, mock.Anything)
}

// OnRunOutput set an expectation on calling executable and printing output to stdout
func (m *Executable) OnRunOutput(path string, args []string, output string) *mock.Call {
	return m.OnRun(path, args).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(internal.RunOptions)
			_, _ = c.Output.Write([]byte(output))
		})
}

// NewExecutable creates a new executable
func (m *Executable) NewExecutable(path string) *internal.Executable {
	return &internal.Executable{Path: path, RunFunc: m.runFunc}
}

// LazyExecutable creates a new lazy executable via function call
func (m *Executable) LazyExecutable(path string) func() *internal.Executable {
	return func() *internal.Executable {
		return m.NewExecutable(path)
	}
}
