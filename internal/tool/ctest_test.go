package tool_test

import (
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCTest_DetectCTest(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("ctest").
		Return(env.NewExecutable("/bin/ctest"))

	ctest := tool.DetectCTest(t.Context(), env)
	assert.NotNil(t, ctest)
	assert.Equal(t, "ctest", ctest.Id())
	assert.True(t, ctest.IsAvailable())

	if executable := ctest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/ctest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestCTest_DetectCTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("ctest").
		Return(nil)

	ctest := tool.DetectCTest(t.Context(), env)
	assert.NotNil(t, ctest)
	assert.Equal(t, "ctest", ctest.Id())
	assert.False(t, ctest.IsAvailable())
	assert.Nil(t, ctest.Executable())

	env.AssertExpectations(t)
}

func TestCTest_RunForProject_NoIntermediateDirectory(t *testing.T) {
	exec := test.NewExecutable(t)

	ctest := tool.NewCTest(exec.LazyExecutable("ctest"))

	desc := internal.ProjectInfo{Directory: "project", IntermediateDirectory: nil}

	err := ctest.RunForProject(t.Context(), desc, []string{})
	require.ErrorIs(t, err, internal.ErrNoIntermediateDirectory)

	exec.AssertExpectations(t)
}

func TestCTest_RunForProject(t *testing.T) {
	exec := test.NewExecutable(t)

	ctest := tool.NewCTest(exec.LazyExecutable("ctest"))

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(nil)

	err := ctest.RunForProject(t.Context(), desc, []string{})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestCTest_RunForProject_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	ctest := tool.NewCTest(exec.LazyExecutable("ctest"))

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("ctest", []string{"--test-dir", "build"}).
		Return(errors.New("failed"))

	err := ctest.RunForProject(t.Context(), desc, []string{})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
