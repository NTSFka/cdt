package test

import (
	"cdt/internal"
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
)

type Environment struct {
	mock.Mock
}

// NewEnvironment create new testing environment
func NewEnvironment(t *testing.T) *Environment {
	env := Environment{}
	env.Test(t)
	return &env
}

func (e *Environment) NewExecutable(path string) *internal.Executable {
	return &internal.Executable{Path: path, Runtime: e}
}

func (e *Environment) Id() string {
	return e.Called().Get(0).(string)
}

func (e *Environment) Start(ctx context.Context) error {
	return e.Called(ctx).Error(0)
}

func (e *Environment) IsRunning(ctx context.Context) bool {
	return e.Called(ctx).Get(0).(bool)
}

func (e *Environment) Stop(ctx context.Context) error {
	return e.Called(ctx).Error(0)
}

func (e *Environment) Cleanup(ctx context.Context) error {
	return e.Called(ctx).Error(0)
}

func (e *Environment) FindExecutable(name string) *internal.Executable {
	result := e.Called(name).Get(0)

	if result == nil {
		return nil
	}

	return result.(*internal.Executable)
}

func (e *Environment) RunExecutable(ctx context.Context, options internal.RunOptions, path string, args []string) error {
	return e.Called(ctx, options, path, args).Error(0)
}

func (e *Environment) OnFindExecutable(name string) *mock.Call {
	return e.On("FindExecutable", name)
}
