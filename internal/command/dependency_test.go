package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func runDependency(ctx context.Context, manager internal.ProjectDependencyManager, args ...string) error {
	return test.RunCommand(ctx, NewDependencyCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				DependencyManager: manager,
			},
		},
	}, args...)
}

func runDependencyTool(ctx context.Context, tool internal.Tool, args ...string) error {
	return test.RunCommand(ctx, NewDependencyCommand(), internal.Context{
		Tools: []internal.Tool{
			tool,
		},
	}, args...)
}

func TestDependency_NotSupported(t *testing.T) {
	err := runDependency(context.Background(), nil, "list")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support dependency management", err.Error())
	}
}

func TestDependency_AddDependencies_Success(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("AddDependencies", mock.Anything, mock.Anything, []string{"dep1", "dep2"}, false).
		Return(nil)

	err := runDependency(context.Background(), manager, "add", "dep1", "dep2")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_AddDependencies_Success_Dev(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("AddDependencies", mock.Anything, mock.Anything, []string{"dep1", "dep2"}, true).
		Return(nil)

	err := runDependency(context.Background(), manager, "add", "--dev", "dep1", "dep2")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_AddDependencies_Failure(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("AddDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).
		Return(errors.New("failed"))

	err := runDependency(context.Background(), manager, "add", "dep1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to add dependencies: failed", err.Error())
	}
	manager.AssertExpectations(t)
}

func TestDependency_RemoveDependencies_Success(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("RemoveDependencies", mock.Anything, mock.Anything, []string{"dep1", "dep2"}, false).
		Return(nil)

	err := runDependency(context.Background(), manager, "remove", "dep1", "dep2")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_RemoveDependencies_Success_Dev(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("RemoveDependencies", mock.Anything, mock.Anything, []string{"dep1", "dep2"}, true).
		Return(nil)

	err := runDependency(context.Background(), manager, "remove", "--dev", "dep1", "dep2")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_UpdateDependencies_Success(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("UpdateDependencies", mock.Anything, mock.Anything, []string{"dep1", "dep2"}).
		Return(nil)

	err := runDependency(context.Background(), manager, "update", "dep1", "dep2")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_UpdateDependencies_Failure(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("UpdateDependencies", mock.Anything, mock.Anything, []string{"dep1"}).
		Return(errors.New("failed"))

	err := runDependency(context.Background(), manager, "update", "dep1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to update dependencies: failed", err.Error())
	}
	manager.AssertExpectations(t)
}

func TestDependency_FetchDependencies_Success(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("FetchDependencies", mock.Anything, mock.Anything, false).
		Return(nil)

	err := runDependency(context.Background(), manager, "fetch")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_FetchDependencies_Success_NoDev(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("FetchDependencies", mock.Anything, mock.Anything, true).
		Return(nil)

	err := runDependency(context.Background(), manager, "fetch", "--no-dev")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_FetchDependencies_Failure(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("FetchDependencies", mock.Anything, mock.Anything, false).
		Return(errors.New("failed"))

	err := runDependency(context.Background(), manager, "fetch")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to fetch dependencies: failed", err.Error())
	}
	manager.AssertExpectations(t)
}

func TestDependency_ListDependencies_Success(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("ListDependencies", mock.Anything, mock.Anything).
		Return(nil)

	err := runDependency(context.Background(), manager, "list")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_ListDependencies_Failure(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("ListDependencies", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runDependency(context.Background(), manager, "list")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to list dependencies: failed", err.Error())
	}
	manager.AssertExpectations(t)
}

func TestDependency_AuditDependencies_Success(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("AuditDependencies", mock.Anything, mock.Anything).
		Return(nil)

	err := runDependency(context.Background(), manager, "audit")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_AuditDependencies_Failure(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("AuditDependencies", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runDependency(context.Background(), manager, "audit")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to audit dependencies: failed", err.Error())
	}
	manager.AssertExpectations(t)
}

func TestDependency_RemoveDependencies_Failure(t *testing.T) {
	manager := test.NewDependencyManager(t)
	manager.On("RemoveDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).
		Return(errors.New("failed"))

	err := runDependency(context.Background(), manager, "remove", "dep1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to remove dependencies: failed", err.Error())
	}
	manager.AssertExpectations(t)
}

type testDependencyTool struct {
	internal.ExecutableTool
	test.DependencyManager
}

func newTestDependencyTool(t *testing.T) *testDependencyTool {
	manager := &testDependencyTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.DependencyManager{},
	}
	manager.Test(t)
	return manager
}

func TestDependency_Tool_Success(t *testing.T) {
	manager := newTestDependencyTool(t)
	manager.On("ListDependencies", mock.Anything, mock.Anything).
		Return(nil)

	err := runDependencyTool(context.Background(), manager, "--tool", "tool1", "list")

	assert.NoError(t, err)
	manager.AssertExpectations(t)
}

func TestDependency_Tool_Failed(t *testing.T) {
	manager := newTestDependencyTool(t)
	manager.On("ListDependencies", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runDependencyTool(context.Background(), manager, "--tool", "tool1", "list")

	if assert.Error(t, err) {
		assert.Equal(t, "failed to list dependencies: failed", err.Error())
	}
	manager.AssertExpectations(t)
}

func TestDependency_Tool_NotFound(t *testing.T) {
	manager := newTestDependencyTool(t)

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
			err := runDependencyTool(context.Background(), manager, append([]string{"--tool", "tool2"}, data...)...)

			if assert.Error(t, err) {
				assert.Equal(t, "tool 'tool2' not found", err.Error())
			}
			manager.AssertExpectations(t)
		})
	}
}

func TestDependency_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
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
			err := runDependencyTool(context.Background(), &linter, append([]string{"--tool", "tool1"}, data...)...)

			if assert.Error(t, err) {
				assert.Equal(t, "tool 'tool1' doesn't support dependency management", err.Error())
			}
		})
	}
}
