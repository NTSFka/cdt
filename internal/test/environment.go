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

func (e *Environment) RunExecutable(path string, args []string) error {
	return e.Called(path, args).Error(0)
}
