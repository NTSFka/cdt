package test

import (
	"cdt/internal"
	"github.com/stretchr/testify/mock"
)

type Environment struct {
	mock.Mock
}

func (e *Environment) MakeExecutable(path string) *internal.Executable {
	return &internal.Executable{Path: path, RunFunc: e.RunExecutable}
}

func (e *Environment) FindExecutable(name string) *internal.Executable {
	result := e.Called(name).Get(0)

	if result == nil {
		return nil
	}

	return result.(*internal.Executable)
}

func (e *Environment) RunExecutable(ctx internal.RunContext, path string, args []string) error {
	return e.Called(ctx, path, args).Error(0)
}

func (e *Environment) OnRun(project internal.Project, path string, args []string) *mock.Call {
	return e.On("RunExecutable", internal.NewRunContext(project.RootDirectory()), path, args)
}

func (e *Environment) OnRunSuccess(project internal.Project, path string, args []string) {
	e.OnRun(project, path, args).Return(nil)
}

func (e *Environment) OnRunError(project internal.Project, path string, args []string, result error) {
	e.OnRun(project, path, args).Return(result)
}
