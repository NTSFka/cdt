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

func TestBuildNotSupported(t *testing.T) {
	err := runBuild(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support building", err.Error())
	}
}

func TestBuildAllSuccess(t *testing.T) {
	builder := test.ProjectBuilder{}
	builder.On("BuildAll", mock.Anything, []string{}).Return(nil)

	err := runBuild(&builder)

	assert.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuildAllFailure(t *testing.T) {
	builder := test.ProjectBuilder{}
	builder.On("BuildAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runBuild(&builder)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	builder.AssertExpectations(t)
}

func TestBuildTargetsSuccess(t *testing.T) {
	builder := test.ProjectBuilder{}
	builder.On("BuildTargets", mock.Anything, []string{"target1", "target2"}, []string{}).Return(nil)

	err := runBuild(&builder, "target1", "target2")

	assert.NoError(t, err)
	builder.AssertExpectations(t)
}

func TestBuildTargetsFailure(t *testing.T) {
	builder := test.ProjectBuilder{}
	builder.On("BuildTargets", mock.Anything, []string{"target1", "target2"}, []string{}).Return(errors.New("failed"))

	err := runBuild(&builder, "target1", "target2")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	builder.AssertExpectations(t)
}
