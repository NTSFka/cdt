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

func runCoverage(
	ctx context.Context,
	collector internal.ProjectCoverageCollector,
	args ...string,
) error {
	return test.RunCommand(ctx, command.NewCoverageCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				CoverageCollector: collector,
			},
		},
	}, args...)
}

func runCoverageTool(ctx context.Context, collector internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewCoverageCommand(), internal.Context{
		Tools: []internal.Tool{
			collector,
		},
	}, args...)
}

func TestCoverage_CannotBeTested(t *testing.T) {
	err := runCoverage(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support coverage collection", err.Error())
}

func TestCoverage_CollectCoverage_Success(t *testing.T) {
	collector := test.NewProjectCoverageCollector(t)
	collector.On("CollectCoverage", mock.Anything, mock.Anything).
		Return(nil)

	err := runCoverage(t.Context(), collector)

	require.NoError(t, err)
	collector.AssertExpectations(t)
}

func TestCoverage_CollectCoverage_Failure(t *testing.T) {
	collector := test.NewProjectCoverageCollector(t)
	collector.On("CollectCoverage", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runCoverage(t.Context(), collector)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())
	collector.AssertExpectations(t)
}

func TestCoverage_CollectCoverage_Pattern_Success(t *testing.T) {
	collector := test.NewProjectCoverageCollector(t)
	collector.On(
		"CollectCoverage",
		mock.Anything,
		mock.MatchedBy(func(options internal.ProjectCoverageCollectorOptions) bool {
			return assert.NotNil(t, options.Pattern) && assert.Equal(t, "pattern", *options.Pattern)
		}),
	).
		Return(nil)

	err := runCoverage(t.Context(), collector, "pattern")

	require.NoError(t, err)
	collector.AssertExpectations(t)
}

type testCoverageCollectorTool struct {
	internal.ExecutableTool
	test.ProjectCoverageCollector
}

func newTestCoverageCollectorTool(t *testing.T) *testCoverageCollectorTool {
	collector := &testCoverageCollectorTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectCoverageCollector{},
	}

	collector.Test(t)

	return collector
}

func TestCoverage_Tool_Success(t *testing.T) {
	collector := newTestCoverageCollectorTool(t)
	collector.On("CollectCoverage", mock.Anything, mock.Anything).
		Return(nil)

	err := runCoverageTool(t.Context(), collector, "--tool", "tool1")

	require.NoError(t, err)
	collector.AssertExpectations(t)
}

func TestCoverage_Tool_Failed(t *testing.T) {
	collector := newTestCoverageCollectorTool(t)
	collector.On("CollectCoverage", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runCoverageTool(t.Context(), collector, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	collector.AssertExpectations(t)
}

func TestCoverage_Tool_NotFound(t *testing.T) {
	collector := newTestCoverageCollectorTool(t)

	err := runCoverageTool(t.Context(), collector, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())

	collector.AssertExpectations(t)
}

func TestCoverage_Tool_NotSupported(t *testing.T) {
	collector := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runCoverageTool(t.Context(), &collector, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support coverage collection", err.Error())
}
