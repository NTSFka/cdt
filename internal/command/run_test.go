package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func runRun(runner internal.ProjectRunner, args ...string) error {
	return test.RunCommand(RunCommand, internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Runner: runner,
			},
		},
	}, args...)
}

func TestRunNotSupported(t *testing.T) {
	err := runRun(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support run of target", err.Error())
	}
}

func TestRunAllSuccess(t *testing.T) {
	runner := test.ProjectRunner{}
	runner.On("RunTarget", mock.Anything, "target1", []string{}).Return(nil)

	err := runRun(&runner, "target1")

	assert.NoError(t, err)
	runner.AssertExpectations(t)
}

func TestRunAllFailure(t *testing.T) {
	runner := test.ProjectRunner{}
	runner.On("RunTarget", mock.Anything, "target1", []string{}).Return(errors.New("failed"))

	err := runRun(&runner, "target1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	runner.AssertExpectations(t)
}
