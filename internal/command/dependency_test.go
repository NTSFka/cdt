package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runDependency(manager internal.ProjectDependencyManager, args ...string) error {
	return test.RunCommand(NewDependencyCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				DependencyManager: manager,
			},
		},
	}, args...)
}

func runDependencyTool(tool internal.Tool, args ...string) error {
	return test.RunCommand(NewDependencyCommand(), internal.Context{
		Project: internal.Project{},
		Tools: []internal.Tool{
			tool,
		},
	}, args...)
}

func TestDependency_NotSupported(t *testing.T) {
	err := runDependency(nil, "list")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support dependency management", err.Error())
	}
}

func TestDependency_AddDependencies_Success(t *testing.T) {
	manager := test.DependencyManager{}
	manager.On("AddDependencies", mock.Anything, []string{"dep1", "dep2"}).Return(nil)

	err := runDependency(&manager, "add", "dep1", "dep2")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_AddDependencies_Failure(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("AddDependencies", mock.Anything, []string{"dep1"}).Return(errors.New("failed"))

	err := runDependency(&configurator, "add", "dep1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to add dependencies: failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestDependency_RemoveDependencies_Success(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("RemoveDependencies", mock.Anything, []string{"dep1", "dep2"}).Return(nil)

	err := runDependency(&configurator, "remove", "dep1", "dep2")

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestDependency_UpdateDependencies_Success(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("UpdateDependencies", mock.Anything, []string{"dep1", "dep2"}).Return(nil)

	err := runDependency(&configurator, "update", "dep1", "dep2")

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestDependency_UpdateDependencies_Failure(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("UpdateDependencies", mock.Anything, []string{"dep1"}).Return(errors.New("failed"))

	err := runDependency(&configurator, "update", "dep1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to update dependencies: failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestDependency_FetchDependencies_Success(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("FetchDependencies", mock.Anything).Return(nil)

	err := runDependency(&configurator, "fetch")

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestDependency_FetchDependencies_Failure(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("FetchDependencies", mock.Anything).Return(errors.New("failed"))

	err := runDependency(&configurator, "fetch")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to fetch dependencies: failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestDependency_ListDependencies_Success(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("ListDependencies", mock.Anything).Return(nil)

	err := runDependency(&configurator, "list")

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestDependency_ListDependencies_Failure(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("ListDependencies", mock.Anything).Return(errors.New("failed"))

	err := runDependency(&configurator, "list")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to list dependencies: failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestDependency_AuditDependencies_Success(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("AuditDependencies", mock.Anything).Return(nil)

	err := runDependency(&configurator, "audit")

	assert.NoError(t, err)
	configurator.AssertExpectations(t)
}

func TestDependency_AuditDependencies_Failure(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("AuditDependencies", mock.Anything).Return(errors.New("failed"))

	err := runDependency(&configurator, "audit")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to audit dependencies: failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestDependency_RemoveDependencies_Failure(t *testing.T) {
	configurator := test.DependencyManager{}
	configurator.On("RemoveDependencies", mock.Anything, []string{"dep1"}).Return(errors.New("failed"))

	err := runDependency(&configurator, "remove", "dep1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to remove dependencies: failed", err.Error())
	}
	configurator.AssertExpectations(t)
}

func TestDependency_Tool_Success(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.DependencyManager
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.DependencyManager{},
	}
	linter.On("ListDependencies", mock.Anything).Return(nil)

	err := runDependencyTool(&linter, "--tool", "tool1", "list")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestDependency_Tool_Failed(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.DependencyManager
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.DependencyManager{},
	}
	linter.On("ListDependencies", mock.Anything).Return(errors.New("failed"))

	err := runDependencyTool(&linter, "--tool", "tool1", "list")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to list dependencies: failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestDependency_Tool_NotFound(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.DependencyManager
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.DependencyManager{},
	}

	dataSet := [][]string{
		{"add", "dep1", "dep2"},
		{"remove", "dep1", "dep2"},
		{"update", "dep1", "dep2"},
		{"fetch"},
		{"list"},
		{"audit"},
	}

	for _, data := range dataSet {
		t.Run(data[0], func(t *testing.T) {
			err := runDependencyTool(&linter, append([]string{"--tool", "tool2"}, data...)...)

			if assert.Error(t, err) {
				assert.Equal(t, "tool 'tool2' not found", err.Error())
			}
			linter.AssertExpectations(t)
		})
	}
}

func TestDependency_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
	}

	dataSet := [][]string{
		{"add", "dep1", "dep2"},
		{"remove", "dep1", "dep2"},
		{"update", "dep1", "dep2"},
		{"fetch"},
		{"list"},
		{"audit"},
	}

	for _, data := range dataSet {
		t.Run(data[0], func(t *testing.T) {
			err := runDependencyTool(&linter, append([]string{"--tool", "tool1"}, data...)...)

			if assert.Error(t, err) {
				assert.Equal(t, "tool 'tool1' doesn't support dependency management", err.Error())
			}
		})
	}
}
