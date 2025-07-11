package test

import (
	"cdt/internal"
	"context"
	"github.com/stretchr/testify/mock"
)

type Environment struct {
	mock.Mock
}

func (e *Environment) MakeExecutable(path string) *internal.Executable {
	return &internal.Executable{Path: path, RunFunc: e.RunExecutable}
}

func (e *Environment) DetectExecutable(path string) func() *internal.Executable {
	return func() *internal.Executable {
		return e.MakeExecutable(path)
	}
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

func (e *Environment) OnRun(project internal.Project, path string, args []string) *mock.Call {
	return e.On("RunExecutable", mock.Anything, internal.RunOptions{Directory: project.RootDirectory()}, path, args)
}

func (e *Environment) OnRunSuccess(project internal.Project, path string, args []string) {
	e.OnRun(project, path, args).Return(nil)
}

func (e *Environment) OnRunError(project internal.Project, path string, args []string, result error) {
	e.OnRun(project, path, args).Return(result)
}
