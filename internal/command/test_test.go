package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func runTest(tester internal.ProjectTester, args ...string) error {
	return test.RunCommand(NewTestCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Tester: tester,
			},
		},
	}, args...)
}

func runTestTool(tester internal.Tool, args ...string) error {
	return test.RunCommand(NewTestCommand(), internal.Context{
		Tools: []internal.Tool{
			tester,
		},
	}, args...)
}

func TestTest_CannotBeTested(t *testing.T) {
	err := runTest(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support testing", err.Error())
	}
}

func TestTest_TestAll_Success(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("TestAll", mock.Anything, []string{}).
		Return(nil)

	err := runTest(tester)

	assert.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_TestAll_Failure(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("TestAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runTest(tester)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	tester.AssertExpectations(t)
}

type testTesterTool struct {
	internal.ExecutableTool
	test.ProjectTester
}

func newTestTesterTool(t *testing.T) *testTesterTool {
	tester := &testTesterTool{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
		test.ProjectTester{},
	}
	tester.Mock.Test(t)
	return tester
}

func TestTest_Tool_Success(t *testing.T) {
	tester := newTestTesterTool(t)
	tester.On("TestAll", mock.Anything, []string{}).
		Return(nil)

	err := runTestTool(tester, "--tool", "tool1")

	assert.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_Tool_Failed(t *testing.T) {
	tester := newTestTesterTool(t)
	tester.On("TestAll", mock.Anything, []string{}).
		Return(errors.New("failed"))

	err := runTestTool(tester, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	tester.AssertExpectations(t)
}

func TestTest_Tool_NotFound(t *testing.T) {
	tester := newTestTesterTool(t)

	err := runTestTool(tester, "--tool", "tool2")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool2' not found", err.Error())
	}
	tester.AssertExpectations(t)
}

func TestTest_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runTestTool(&linter, "--tool", "tool1")

	if assert.Error(t, err) {
		assert.Equal(t, "tool 'tool1' doesn't support testing", err.Error())
	}
}

func TestTest_TestTargets_Success(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("Test", mock.Anything, "pattern", []string{}).
		Return(nil)

	err := runTest(tester, "pattern")

	assert.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_TestTargets_Failure(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("Test", mock.Anything, "pattern", []string{}).
		Return(errors.New("failed"))

	err := runTest(tester, "pattern")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	tester.AssertExpectations(t)
}
