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

func runBuild(ctx context.Context, builder internal.ProjectBuilder, args ...string) error {
	return test.RunCommand(ctx, command.NewBuildCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Builder: builder,
			},
		},
	}, args...)
}

func runBuildTool(ctx context.Context, builder internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewBuildCommand(), internal.Context{
		Tools: []internal.Tool{
			builder,
		},
	}, args...)
}

func TestBuild_NotSupported(t *testing.T) {
	err := runBuild(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support building", err.Error())
}

func TestBuild_BuildAll_Success(t *testing.T) {
	builder := test.NewProjectBuilder(t)
	builder.On("BuildAll", mock.Anything, mock.Anything).
		Return(nil)

	err := runBuild(t.Context(), builder)

	require.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuild_BuildAll_Failure(t *testing.T) {
	builder := test.NewProjectBuilder(t)
	builder.On("BuildAll", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runBuild(t.Context(), builder)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	builder.AssertExpectations(t)
}

type testBuilderTool struct {
	internal.ExecutableTool
	test.ProjectBuilder
}

func newTestBuilderTool(t *testing.T) *testBuilderTool {
	builder := testBuilderTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectBuilder{},
	}
	builder.Test(t)

	return &builder
}

func TestBuild_Tool_Success(t *testing.T) {
	builder := newTestBuilderTool(t)
	builder.On("BuildAll", mock.Anything, mock.Anything).
		Return(nil)

	err := runBuildTool(t.Context(), builder, "--tool", "tool1")

	require.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuild_Tool_Failed(t *testing.T) {
	builder := newTestBuilderTool(t)
	builder.On("BuildAll", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runBuildTool(t.Context(), builder, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	builder.AssertExpectations(t)
}

func TestBuild_Tool_NotFound(t *testing.T) {
	builder := newTestBuilderTool(t)

	err := runBuildTool(t.Context(), builder, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())

	builder.AssertExpectations(t)
}

func TestBuild_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runBuildTool(t.Context(), &linter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support building", err.Error())
}

func TestBuild_BuildTargets_Success(t *testing.T) {
	builder := test.NewProjectBuilder(t)

	builder.On("BuildTargets", mock.Anything, mock.Anything, []string{"target1", "target2"}).
		Return(nil)

	err := runBuild(t.Context(), builder, "target1", "target2")

	require.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuild_BuildTargets_Failure(t *testing.T) {
	builder := test.NewProjectBuilder(t)

	builder.On("BuildTargets", mock.Anything, mock.Anything, []string{"target1", "target2"}).
		Return(errors.New("failed"))

	err := runBuild(t.Context(), builder, "target1", "target2")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	builder.AssertExpectations(t)
}
