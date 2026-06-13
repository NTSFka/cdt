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

	coverage := tool.DetectPythonCoverage(t.Context(), internal.DetectOptions{Environment: env})
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

	coverage := tool.DetectPythonCoverage(t.Context(), internal.DetectOptions{Environment: env})
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

	coverage := tool.DetectPythonCoverage(t.Context(), internal.DetectOptions{
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
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), func() string {
		return "my-test"
	})

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "my-test"}).
		Return(nil)

	exec.OnRun("test", []string{"report"}).
		Return(nil)

	err := coverage.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPythonCoverage_PythonCoverage_CollectCoverage_FailCollection(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), func() string {
		return "my-test"
	})

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "my-test"}).
		Return(errors.New("failed"))

	err := coverage.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestPythonCoverage_PythonCoverage_CollectCoverage_Pattern(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), nil)

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "unittest", "tests/*"}).
		Return(nil)

	exec.OnRun("test", []string{"report"}).
		Return(nil)

	err := coverage.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: info,
			Pattern:     internal.StrPtr("tests/*"),
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPythonCoverage_PythonCoverage_CollectCoverage_Pattern_FailedCollection(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), nil)

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"run", "-m", "unittest", "tests/*"}).
		Return(errors.New("failed"))

	err := coverage.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: info,
			Pattern:     internal.StrPtr("tests/*"),
		},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestPythonCoverage_CollectCoverage_FormatUnsupported(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), nil)

	exec.OnRun("test", []string{"run", "-m", "unittest", "tests/*"}).
		Return(nil)

	info := internal.ProjectInfo{Directory: "."}
	data := []internal.CoverageReportFormat{
		"test",
		internal.CoverageReportFormatCobertura,
		internal.CoverageReportFormatCrap4j,
	}

	for _, format := range data {
		t.Run(string(format), func(t *testing.T) {
			err := coverage.CollectCoverage(
				t.Context(),
				internal.ProjectCoverageCollectorOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.CoverageReportFormat]{
						Format: format,
					},
					Pattern: internal.StrPtr("tests/*"),
				},
			)
			require.EqualError(t, err, "unsupported coverage report format: "+string(format))

			exec.AssertExpectations(t)
		})
	}
}

func TestPythonCoverage_CollectCoverage_FormatSupported(t *testing.T) {
	exec := test.NewExecutable(t)

	coverage := tool.NewPythonCoverage(exec.LazyExecutable("test"), nil)

	exec.OnRun("test", []string{"run", "-m", "unittest", "tests/*"}).
		Return(nil)

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
			exec.OnRun("test", values.Args).
				Return(nil)

			err := coverage.CollectCoverage(
				t.Context(),
				internal.ProjectCoverageCollectorOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.CoverageReportFormat]{
						Format:   values.Format,
						Filename: values.Filename,
					},
					Pattern: internal.StrPtr("tests/*"),
				},
			)
			require.NoError(t, err)

			exec.AssertExpectations(t)
		})
	}
}
