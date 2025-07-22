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
	return test.RunCommand(NewRunCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Runner: runner,
			},
		},
	}, args...)
}

func runRunTool(runner internal.Tool, args ...string) error {
	return test.RunCommand(NewRunCommand(), internal.Context{
		Project: internal.Project{},
		Tools: []internal.Tool{
			runner,
		},
	}, args...)
}

func TestRun_NotSupported(t *testing.T) {
	err := runRun(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support run of target", err.Error())
	}
}

func TestRun_Target_NoTarget(t *testing.T) {
	runner := test.ProjectRunner{}

	err := runRun(&runner)

	assert.EqualError(t, err, "target is required")
	runner.AssertExpectations(t)
}

func TestRun_Target_Success(t *testing.T) {
	runner := test.ProjectRunner{}
	runner.On("RunTarget", mock.Anything, "target1", []string{}).Return(nil)

	err := runRun(&runner, "target1")

	assert.NoError(t, err)
	runner.AssertExpectations(t)
}

func TestRun_Target_Failure(t *testing.T) {
	runner := test.ProjectRunner{}
	runner.On("RunTarget", mock.Anything, "target1", []string{}).Return(errors.New("failed"))

	err := runRun(&runner, "target1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	runner.AssertExpectations(t)
}

func TestRun_Tool_Success(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectRunner
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectRunner{},
	}
	linter.On("RunTarget", mock.Anything, "target1", []string{}).Return(nil)

	err := runRunTool(&linter, "--tool", "tool1", "target1")

	assert.NoError(t, err)
	linter.AssertExpectations(t)
}

func TestRun_Tool_Failed(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectRunner
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectRunner{},
	}
	linter.On("RunTarget", mock.Anything, "target1", []string{}).Return(errors.New("failed"))

	err := runRunTool(&linter, "--tool", "tool1", "target1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestRun_Tool_NotFound(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
		test.ProjectRunner
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
		test.ProjectRunner{},
	}

	err := runRunTool(&linter, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	linter.AssertExpectations(t)
}

func TestRun_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", nil),
	}

	err := runRunTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support run of target", err.Error())
	}
}
