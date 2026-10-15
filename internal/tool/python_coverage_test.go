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

func TestPythonCoverage_PythonCoverage_CollectCoverage(t *testing.T) {
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

func TestPythonCoverage_PythonCoverage_CollectCoverage_FailCollection(t *testing.T) {
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

func TestPythonCoverage_PythonCoverage_CollectCoverage_Pattern(t *testing.T) {
	test.InvokeCollectCoverage(
		t,
		pythonCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     internal.StrPtr("tests/*"),
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

func TestPythonCoverage_PythonCoverage_CollectCoverage_Pattern_FailedCollection(t *testing.T) {
	test.InvokeCollectCoverage(
		t,
		pythonCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     internal.StrPtr("tests/*"),
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

func TestPythonCoverage_CollectCoverage_FormatUnsupported(t *testing.T) {
	info := internal.ProjectInfo{Directory: "."}
	data := []internal.CoverageReportFormat{
		"test",
		internal.CoverageReportFormatCobertura,
		internal.CoverageReportFormatCrap4j,
	}

	for _, format := range data {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeCollectCoverage(
				t,
				pythonCoverageCollector,
				internal.ProjectCoverageCollectorOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.CoverageReportFormat]{
						Format: format,
					},
					Pattern: internal.StrPtr("tests/*"),
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRun(executablePath, []string{"run", "-m", "unittest", "tests/*"}).
						Return(nil)
				},
				func(err error) {
					require.EqualError(
						t,
						err,
						"unsupported coverage report format: "+string(format),
					)
				},
			)
		})
	}
}

func TestPythonCoverage_CollectCoverage_FormatSupported(t *testing.T) {
	info := internal.ProjectInfo{Directory: "."}
	data := []struct {
		Format   internal.CoverageReportFormat
		Filename *string
		Args     []string
	}{
		{
			Format: internal.CoverageReportFormatRaw,
			Args:   []string{"report"},
		},
		{
			Format:   internal.CoverageReportFormatHtml,
			Filename: internal.StrPtr("coverage.html"),
			Args:     []string{"html", "coverage.html"},
		},
		{
			Format:   internal.CoverageReportFormatXml,
			Filename: internal.StrPtr("coverage.xml"),
			Args:     []string{"xml", "coverage.xml"},
		},
		{
			Format:   internal.CoverageReportFormatJson,
			Filename: internal.StrPtr("coverage.json"),
			Args:     []string{"json", "coverage.json"},
		},
		{
			Format:   internal.CoverageReportFormatLcov,
			Filename: internal.StrPtr("coverage.lcov"),
			Args:     []string{"lcov", "coverage.lcov"},
		},
	}

	for _, values := range data {
		t.Run(string(values.Format), func(t *testing.T) {
			test.InvokeCollectCoverage(
				t,
				pythonCoverageCollector,
				internal.ProjectCoverageCollectorOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.CoverageReportFormat]{
						Format:   values.Format,
						Filename: values.Filename,
					},
					Pattern: internal.StrPtr("tests/*"),
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRun(executablePath, []string{"run", "-m", "unittest", "tests/*"}).
						Return(nil)

					exec.OnRun(executablePath, values.Args).
						Return(nil)
				},
				func(err error) {
					require.NoError(t, err)
				},
			)
		})
	}
}

func pythonCoverageDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPythonCoverage(ctx, options)
}

func pythonCoverageCollector(
	executable func() (*internal.Executable, error),
) internal.ProjectCoverageCollector {
	return tool.NewPythonCoverage(executable, nil)
}
