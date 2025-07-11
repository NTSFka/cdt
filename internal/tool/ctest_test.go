package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCTest_DetectCTest(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("ctest").
		Return(environment.MakeExecutable("/bin/ctest"))

	tool := DetectCTest(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "ctest", tool.Id())
	assert.True(t, tool.IsAvailable())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/ctest", *path)
	}

	environment.AssertExpectations(t)
}

func TestCTest_DetectCTest_NotFound(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("ctest").
		Return(nil)

	tool := DetectCTest(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "ctest", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestCTest_Run(t *testing.T) {
	runMock := test.Executable{}

	tool := NewCTest(runMock.LazyExecutable("ctest"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(nil)

	err := tool.RunForProject(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestCTest_Run_Failed(t *testing.T) {
	runMock := test.Executable{}

	tool := NewCTest(runMock.LazyExecutable("ctest"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	runMock.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(errors.New("failed"))

	err := tool.RunForProject(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}
