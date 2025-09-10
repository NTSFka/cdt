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

func runRun(ctx context.Context, runner internal.ProjectRunner, args ...string) error {
	return test.RunCommand(ctx, NewRunCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Runner: runner,
			},
		},
	}, args...)
}

func runRunTool(ctx context.Context, runner internal.Tool, args ...string) error {
	return test.RunCommand(ctx, NewRunCommand(), internal.Context{
		Tools: []internal.Tool{
			runner,
		},
	}, args...)
}

func TestRun_NotSupported(t *testing.T) {
	err := runRun(context.Background(), nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support run of target", err.Error())
	}
}

func TestRun_Target_NoTarget(t *testing.T) {
	runner := test.NewProjectRunner(t)

	err := runRun(context.Background(), runner)

	assert.EqualError(t, err, "target is required")
	runner.AssertExpectations(t)
}

func TestRun_Target_Success(t *testing.T) {
	runner := test.NewProjectRunner(t)
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1", []string{}).
		Return(nil)

	err := runRun(context.Background(), runner, "target1")

	assert.NoError(t, err)
	runner.AssertExpectations(t)
}

func TestRun_Target_Failure(t *testing.T) {
	runner := test.NewProjectRunner(t)
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1", []string{}).
		Return(errors.New("failed"))

	err := runRun(context.Background(), runner, "target1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	runner.AssertExpectations(t)
}

type testRunnerTool struct {
	internal.ExecutableTool
	test.ProjectRunner
}

func newTestRunnerTool(t *testing.T) *testRunnerTool {
	runner := &testRunnerTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectRunner{},
	}
	runner.Test(t)
	return runner
}

func TestRun_Tool_Success(t *testing.T) {
	runner := newTestRunnerTool(t)
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1", []string{}).
		Return(nil)

	err := runRunTool(context.Background(), runner, "--tool", "tool1", "target1")

	assert.NoError(t, err)
	runner.AssertExpectations(t)
}

func TestRun_Tool_Failed(t *testing.T) {
	runner := newTestRunnerTool(t)
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1", []string{}).
		Return(errors.New("failed"))

	err := runRunTool(context.Background(), runner, "--tool", "tool1", "target1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	runner.AssertExpectations(t)
}

func TestRun_Tool_NotFound(t *testing.T) {
	runner := newTestRunnerTool(t)

	err := runRunTool(context.Background(), runner, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	runner.AssertExpectations(t)
}

func TestRun_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runRunTool(context.Background(), &linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support run of target", err.Error())
	}
}
