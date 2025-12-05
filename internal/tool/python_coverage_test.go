package tool_test

import (
	"cdt/internal"
	"errors"
	"testing"

	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPythonCoverage_DetectPythonCoverage(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("coverage").
		Return(env.NewExecutable("/bin/coverage"), nil)

	coverage := tool.DetectPythonCoverage(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, coverage)
	assert.Equal(t, "python-coverage", coverage.Id())
	assert.True(t, coverage.IsAvailable())

	if executable := coverage.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/coverage", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPythonCoverage_DetectPythonCoverage_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("coverage").
		Return(nil, nil)

	coverage := tool.DetectPythonCoverage(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, coverage)
	assert.Equal(t, "python-coverage", coverage.Id())
	assert.False(t, coverage.IsAvailable())
	assert.Nil(t, coverage.Executable())

	env.AssertExpectations(t)
}

func TestPythonCoverage_DetectPythonCoverage_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("coverage2").
		Return(env.NewExecutable("/bin/coverage"), nil)

	coverage := tool.DetectPythonCoverage(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{tool.IdPythonCoverage: "coverage2"},
	})
	assert.NotNil(t, coverage)
	assert.Equal(t, "python-coverage", coverage.Id())
	assert.True(t, coverage.IsAvailable())

	if executable := coverage.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/coverage", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPythonCoverage_DetectPythonCoverage_DetectRunModule_Default(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("coverage").
		Return(env.NewExecutable("/bin/coverage"), nil)

	coverage := tool.DetectPythonCoverage(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, coverage)
	assert.Equal(t, "python-coverage", coverage.Id())
	assert.True(t, coverage.IsAvailable())

	if executable := coverage.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/coverage", executable.Path)
	}

	env.OnFindExecutable("pytest").
		Return(nil, nil)

	assert.Equal(t, "unittest", coverage.DetectRunModule())

	env.AssertExpectations(t)
}

func TestPythonCoverage_DetectPythonCoverage_DetectRunModule_Pytest(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("coverage").
		Return(env.NewExecutable("/bin/coverage"), nil)

	coverage := tool.DetectPythonCoverage(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, coverage)
	assert.Equal(t, "python-coverage", coverage.Id())
	assert.True(t, coverage.IsAvailable())

	if executable := coverage.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/coverage", executable.Path)
	}

	env.OnFindExecutable("pytest").
		Return(env.NewExecutable("/bin/pytest"), nil)

	assert.Equal(t, "pytest", coverage.DetectRunModule())

	env.AssertExpectations(t)
}

func TestPythonCoverage_PythonCoverage_CollectCoverageAll(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), func() string {
		return "my-test"
	})

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "my-test"}).
		Return(nil)

	exec.OnRun("test", []string{"report"}).
		Return(nil)

	err := coverage.CollectCoverageAll(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPythonCoverage_PythonCoverage_CollectCoverageAll_FailCollection(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), func() string {
		return "my-test"
	})

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "my-test"}).
		Return(errors.New("failed"))

	err := coverage.CollectCoverageAll(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestPythonCoverage_PythonCoverage_CollectCoveragePattern(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), nil)

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "unittest", "tests/*"}).
		Return(nil)

	exec.OnRun("test", []string{"report"}).
		Return(nil)

	err := coverage.CollectCoveragePattern(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
		"tests/*",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPythonCoverage_PythonCoverage_CollectCoveragePattern_FailedCollection(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), nil)

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "unittest", "tests/*"}).
		Return(errors.New("failed"))

	err := coverage.CollectCoveragePattern(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
		"tests/*",
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
