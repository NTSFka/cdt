package command_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/command"
	"cdt/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func runRun(ctx context.Context, runner internal.ProjectRunner, args ...string) error {
	return test.RunCommand(ctx, command.NewRunCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Runner: runner,
			},
		},
	}, args...)
}

func runRunTool(ctx context.Context, runner internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewRunCommand(), internal.Context{
		Tools: []internal.Tool{
			runner,
		},
	}, args...)
}

func TestRun_NotSupported(t *testing.T) {
	err := runRun(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support run of target", err.Error())
}

func TestRun_Target_NoTarget(t *testing.T) {
	runner := test.NewProjectRunner(t)

	err := runRun(t.Context(), runner)

	require.EqualError(t, err, "target is required")
	runner.AssertExpectations(t)
}

func TestRun_Target_Success(t *testing.T) {
	runner := test.NewProjectRunner(t)
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1").
		Return(nil)

	err := runRun(t.Context(), runner, "target1")

	require.NoError(t, err)
	runner.AssertExpectations(t)
}

func TestRun_Target_Failure(t *testing.T) {
	runner := test.NewProjectRunner(t)
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1").
		Return(errors.New("failed"))

	err := runRun(t.Context(), runner, "target1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

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
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1").
		Return(nil)

	err := runRunTool(t.Context(), runner, "--tool", "tool1", "target1")

	require.NoError(t, err)
	runner.AssertExpectations(t)
}

func TestRun_Tool_Failed(t *testing.T) {
	runner := newTestRunnerTool(t)
	runner.On("RunTarget", mock.Anything, mock.Anything, "target1").
		Return(errors.New("failed"))

	err := runRunTool(t.Context(), runner, "--tool", "tool1", "target1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	runner.AssertExpectations(t)
}

func TestRun_Tool_NotFound(t *testing.T) {
	runner := newTestRunnerTool(t)

	err := runRunTool(t.Context(), runner, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())

	runner.AssertExpectations(t)
}

func TestRun_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runRunTool(t.Context(), &linter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support run of target", err.Error())
}
