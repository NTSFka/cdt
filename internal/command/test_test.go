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

func TestTest_RunTests_Success(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("RunTests", mock.Anything, mock.Anything).
		Return(nil)

	err := runTest(t.Context(), tester)

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_RunTests_CustomOutput_Raw(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("RunTests", mock.Anything, mock.MatchedBy(func(options internal.ProjectTesterOptions) bool {
		return assert.Equal(t, internal.TestsReportFormatRaw, options.Output.Format) &&
			assert.Nil(t, options.Output.Filename)
	})).
		Return(nil)

	err := runTest(t.Context(), tester, "--output", "raw")

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_RunTests_CustomOutput_Json(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("RunTests", mock.Anything, mock.MatchedBy(func(options internal.ProjectTesterOptions) bool {
		return assert.Equal(t, internal.TestsReportFormatJson, options.Output.Format) &&
			assert.NotNil(t, options.Output.Filename) &&
			assert.Equal(t, "tests.json", *options.Output.Filename)
	})).
		Return(nil)

	err := runTest(t.Context(), tester, "-o", "json:tests.json")

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_RunTests_Failure(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On("RunTests", mock.Anything, mock.Anything).
		Return(errors.New("failed"))

	err := runTest(t.Context(), tester)

	require.Error(t, err)
	assert.Equal(t, "failed", err.Error())
	tester.AssertExpectations(t)
}

func TestTest_RunTests_Pattern_Success(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On(
		"RunTests",
		mock.Anything,
		mock.MatchedBy(func(options internal.ProjectTesterOptions) bool {
			return assert.NotNil(t, options.Pattern) && assert.Equal(t, "pattern", *options.Pattern)
		}),
	).
		Return(nil)

	err := runTest(t.Context(), tester, "pattern")

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_RunTests_Pattern_Failure(t *testing.T) {
	tester := test.NewProjectTester(t)
	tester.On(
		"RunTests",
		mock.Anything,
		mock.MatchedBy(func(options internal.ProjectTesterOptions) bool {
			return assert.NotNil(t, options.Pattern) && assert.Equal(t, "pattern", *options.Pattern)
		}),
	).
		Return(errors.New("failed"))

	err := runTest(t.Context(), tester, "pattern")

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
	tester.On("RunTests", mock.Anything, mock.Anything).
		Return(nil)

	err := runTestTool(t.Context(), tester, "--tool", "tool1")

	require.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTest_Tool_Failed(t *testing.T) {
	tester := newTestTesterTool(t)
	tester.On("RunTests", mock.Anything, mock.Anything).
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
