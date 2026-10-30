package tool_test

import (
	"context"
	"errors"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPythonCoverage_DetectPythonCoverage(t *testing.T) {
	coverage := test.RunDetectToolLastFound(t, pythonCoverageDetect, []string{"coverage"})

	assert.Equal(t, tool.IdPythonCoverage, coverage.Id())
}

func TestPythonCoverage_DetectPythonCoverage_NotFound(t *testing.T) {
	coverage := test.RunDetectToolNotFound(t, pythonCoverageDetect, []string{"coverage"})

	assert.Equal(t, tool.IdPythonCoverage, coverage.Id())
}

func TestPythonCoverage_DetectPythonCoverage_Config(t *testing.T) {
	coverage := test.RunDetectToolConfig(t, pythonCoverageDetect, tool.IdPythonCoverage)

	assert.Equal(t, tool.IdPythonCoverage, coverage.Id())
}

func TestPythonCoverage_DetectPythonCoverage_DetectRunModule_Default(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("coverage").
		Return(env.NewExecutable("/bin/coverage"), nil)

	coverage := tool.DetectPythonCoverage(t.Context(), internal.DetectOptions{Environment: env})
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

	coverage := tool.DetectPythonCoverage(t.Context(), internal.DetectOptions{Environment: env})
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

func TestPythonCoverage_CollectCoverage(t *testing.T) {
	test.InvokeCollectCoverage(
		t,
		pythonCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{"run", "-m", "unittest"}).
				Return(nil)

			exec.OnRun(executablePath, []string{"report"}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestPythonCoverage_CollectCoverage_FailCollection(t *testing.T) {
	test.InvokeCollectCoverage(
		t,
		pythonCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{"run", "-m", "unittest"}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func TestPythonCoverage_CollectCoverage_Pattern(t *testing.T) {
	test.InvokeCollectCoverage(
		t,
		pythonCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     new("tests/*"),
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{"run", "-m", "unittest", "tests/*"}).
				Return(nil)

			exec.OnRun(executablePath, []string{"report"}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestPythonCoverage_CollectCoverage_Pattern_FailedCollection(t *testing.T) {
	test.InvokeCollectCoverage(
		t,
		pythonCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     new("tests/*"),
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{"run", "-m", "unittest", "tests/*"}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func pythonCoverageDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPythonCoverage(ctx, options)
}

func pythonCoverageCollector(
	executable func() (*internal.Executable, error),
) internal.ProjectCoverageCollector {
	return tool.NewPythonCoverage(executable, nil)
}
