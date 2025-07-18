package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

func TestTestCannotBeTested(t *testing.T) {
	err := runTest(nil)

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support testing", err.Error())
	}
}

func TestTestAllSuccess(t *testing.T) {
	tester := test.ProjectTester{}
	tester.On("TestAll", mock.Anything, []string{}).Return(nil)

	err := runTest(&tester)

	assert.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTestAllFailure(t *testing.T) {
	tester := test.ProjectTester{}
	tester.On("TestAll", mock.Anything, []string{}).Return(errors.New("failed"))

	err := runTest(&tester)

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	tester.AssertExpectations(t)
}

func TestTestTargetsSuccess(t *testing.T) {
	tester := test.ProjectTester{}
	tester.On("Test", mock.Anything, "pattern", []string{}).Return(nil)

	err := runTest(&tester, "pattern")

	assert.NoError(t, err)
	tester.AssertExpectations(t)
}

func TestTestTargetsFailure(t *testing.T) {
	tester := test.ProjectTester{}
	tester.On("Test", mock.Anything, "pattern", []string{}).Return(errors.New("failed"))

	err := runTest(&tester, "pattern")

	if assert.Error(t, err) {
		assert.Equal(t, "failed", err.Error())
	}
	tester.AssertExpectations(t)
}
