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

func runTest(ctx context.Context, tester internal.ProjectTester, args ...string) error {
	return test.RunCommand(ctx, command.NewTestCommand(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Tester: tester,
			},
		},
	}, args...)
}

func runTestTool(ctx context.Context, tester internal.Tool, args ...string) error {
	return test.RunCommand(ctx, command.NewTestCommand(), internal.Context{
		Tools: []internal.Tool{
			tester,
		},
	}, args...)
}

func TestTest_CannotBeTested(t *testing.T) {
	err := runTest(t.Context(), nil)

	require.Error(t, err)
	assert.Equal(t, "project doesn't support testing", err.Error())
}

func TestTest_TestAll_Success(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("TestAll", mock.Anything, mock.Anything).
		Return(nil)

	err := runTest(t.Context(), tester)

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_TestAll_Failure(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("TestAll", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runTest(t.Context(), tester)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())
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

	tester.Test(t)

	return tester
}

func TestTest_Tool_Success(t *testing.T) {
	tester := newTestTesterTool(t)
	tester.On("TestAll", mock.Anything, mock.Anything).
		Return(nil)

	err := runTestTool(t.Context(), tester, "--tool", "tool1")

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_Tool_Failed(t *testing.T) {
	tester := newTestTesterTool(t)
	tester.On("TestAll", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runTestTool(t.Context(), tester, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	tester.AssertExpectations(t)
}

func TestTest_Tool_NotFound(t *testing.T) {
	tester := newTestTesterTool(t)

	err := runTestTool(t.Context(), tester, "--tool", "tool2")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool2' not found", err.Error())

	tester.AssertExpectations(t)
}

func TestTest_Tool_NotSupported(t *testing.T) {
	linter := struct {
		internal.ExecutableTool
	}{
		internal.MakeExecutableTool("tool1", "", "", internal.Tags{}, nil),
	}

	err := runTestTool(t.Context(), &linter, "--tool", "tool1")

	require.Error(t, err)
	assert.Equal(t, "tool 'tool1' doesn't support testing", err.Error())
}

func TestTest_TestTargets_Success(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("TestPattern", mock.Anything, mock.Anything, "pattern").
		Return(nil)

	err := runTest(t.Context(), tester, "pattern")

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_TestTargets_Failure(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("TestPattern", mock.Anything, mock.Anything, "pattern").
		Return(errors.New("failed"))

	err := runTest(t.Context(), tester, "pattern")

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())

	tester.AssertExpectations(t)
}
