package tool

import (
	"cdt/internal"
	"github.com/stretchr/testify/mock"
	"testing"
)

type testEnvironment struct {
	mock.Mock
}

func (e *testEnvironment) makeExecutable(path string) *internal.Executable {
	return &internal.Executable{Path: path, RunFunc: e.RunExecutable}
}

func (e *testEnvironment) FindExecutable(name string) *internal.Executable {
	result := e.Called(name).Get(0)

	if result == nil {
		return nil
	}

	return result.(*internal.Executable)
}

func (e *testEnvironment) RunExecutable(path string, args []string) error {
	return e.Called(path, args).Error(0)
}

// Check if a tool exists in the given environment
func checkTool(t *testing.T, environment internal.Environment, toolName string) {
	if executable := environment.FindExecutable(toolName); executable == nil {
		t.Skipf("unable to find tool: %v", toolName)
	}
}
