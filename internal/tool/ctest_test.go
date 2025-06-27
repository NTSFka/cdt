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
	environment.On("FindExecutable", "ctest").Return(environment.MakeExecutable("/bin/ctest"))

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
	environment.On("FindExecutable", "ctest").Return(nil)

	tool := DetectCTest(&environment)
	assert.NotNil(t, tool)
	assert.Equal(t, "ctest", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestCTest_Run(t *testing.T) {
	environment := test.Environment{}

	tool := NewCTest(environment.MakeExecutable("ctest"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunSuccess(p, "ctest", []string{"--test-dir", "build"})

	err := tool.Run(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestCTest_Run_Failed(t *testing.T) {
	environment := test.Environment{}

	tool := NewCTest(environment.MakeExecutable("ctest"))

	p := internal.MakeProject("project", "build", nil, internal.Workflow{})

	environment.OnRunError(p, "ctest", []string{"--test-dir", "build"}, errors.New("failed"))

	err := tool.Run(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}
