package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runBuild(builder internal.ProjectBuilder, args ...string) error {
	return test.RunCommand(NewBuildCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Builder: builder,
			},
		},
	}, args...)
}

func runBuildTool(builder internal.Tool, args ...string) error {
	return test.RunCommand(NewBuildCommand(), internal.Context{
		Project: internal.Project{},
		Tools: []internal.Tool{
			builder,
		},
	}, args...)
}

func TestBuild_NotSupported(t *testing.T) {
	err := runBuild(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support building", err.Error())
	}
}

func TestBuild_BuildAll_Success(t *testing.T) {
	builder := test.NewProjectBuilder(t)
	builder.On("BuildAll", mock.Anything, []string{}).
		Return(nil)

	err := runBuild(builder)

	assert.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuild_BuildAll_Failure(t *testing.T) {
	builder := test.NewProjectBuilder(t)
	builder.On("BuildAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runBuild(builder)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	builder.AssertExpectations(t)
}

type testBuilderTool struct {
	internal.ExecutableTool
	test.ProjectBuilder
}

func newTestBuilderTool(t *testing.T) *testBuilderTool {
	builder := testBuilderTool{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectBuilder{},
	}
	builder.Test(t)
	return &builder
}

func TestBuild_Tool_Success(t *testing.T) {
	builder := newTestBuilderTool(t)
	builder.On("BuildAll", mock.Anything, []string{}).
		Return(nil)

	err := runBuildTool(builder, "--tool", "tool1")

	assert.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuild_Tool_Failed(t *testing.T) {
	builder := newTestBuilderTool(t)
	builder.On("BuildAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runBuildTool(builder, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	builder.AssertExpectations(t)
}

func TestBuild_Tool_NotFound(t *testing.T) {
	builder := newTestBuilderTool(t)

	err := runBuildTool(builder, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	builder.AssertExpectations(t)
}

func TestBuild_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
	}

	err := runBuildTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support building", err.Error())
	}
}

func TestBuild_BuildTargets_Success(t *testing.T) {
	builder := test.NewProjectBuilder(t)

	builder.On("BuildTargets", mock.Anything, []string{"target1", "target2"}, []string{}).
		Return(nil)

	err := runBuild(builder, "target1", "target2")

	assert.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuild_BuildTargets_Failure(t *testing.T) {
	builder := test.NewProjectBuilder(t)

	builder.On("BuildTargets", mock.Anything, []string{"target1", "target2"}, []string{}).
		Return(errors.New("failed"))

	err := runBuild(builder, "target1", "target2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	builder.AssertExpectations(t)
}
