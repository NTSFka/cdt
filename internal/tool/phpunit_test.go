package tool_test

import (
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPHPUnit_DetectPHPUnit_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpunit")).
		Return(env.NewExecutable("/bin/phpunit"), nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpunit")).
		Return(nil, nil)

	// System installation
	env.OnFindExecutable("phpunit").
		Return(env.NewExecutable("/bin/phpunit"), nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable(filepath.Join("vendor", "bin", "phpunit")).
		Return(nil, nil)
	env.OnFindExecutable("phpunit").
		Return(nil, nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.False(t, phpUnit.IsAvailable())
	assert.Nil(t, phpUnit.Executable())

	env.AssertExpectations(t)
}

func TestPHPUnit_DetectPHPUnit_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("phpunit-12").
		Return(env.NewExecutable("/bin/phpunit"), nil)

	phpUnit := tool.DetectPHPUnit(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"phpunit": "phpunit-12"},
	})
	assert.NotNil(t, phpUnit)
	assert.Equal(t, "phpunit", phpUnit.Id())
	assert.True(t, phpUnit.IsAvailable())

	if executable := phpUnit.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/phpunit", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_RunTests(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := phpUnit.RunTests(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_RunTests_Pattern(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := phpUnit.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: info, Pattern: internal.StrPtr("tests/*")},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_RunTests_Pattern_FormatUnsupported(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	data := []internal.TestsReportFormat{
		"test",
		internal.TestsReportFormatJson,
		internal.TestsReportFormatCtrf,
	}

	for _, format := range data {
		t.Run(string(format), func(t *testing.T) {
			err := phpUnit.RunTests(
				t.Context(),
				internal.ProjectTesterOptions{
					ProjectInfo: info,
					Output:      internal.OutputOptions[internal.TestsReportFormat]{Format: format},
					Pattern:     internal.StrPtr("test1"),
				},
			)
			require.EqualError(t, err, "unsupported report format: "+string(format))

			exec.AssertExpectations(t)
		})
	}
}

func TestPHPUnit_RunTests_Pattern_FormatRawEvents(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	reportFilename := "phpunit-raw-events.txt"

	exec.OnRun("test", []string{"test1", "--log-events-text", reportFilename}).
		Return(nil)

	err := phpUnit.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatRawEvents,
				Filename: &reportFilename,
			},
			Pattern: internal.StrPtr("test1"),
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_RunTests_Pattern_FormatJUnit(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	reportFilename := "phpunit-junit.xml"

	exec.OnRun("test", []string{"test1", "--log-junit", reportFilename}).
		Return(nil)

	err := phpUnit.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatJUnit,
				Filename: &reportFilename,
			},
			Pattern: internal.StrPtr("test1"),
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_RunTests_Pattern_FormatTeamCity(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	reportFilename := "phpunit-teamcity.xml"

	exec.OnRun("test", []string{"test1", "--log-teamcity", reportFilename}).
		Return(nil)

	err := phpUnit.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatTeamCity,
				Filename: &reportFilename,
			},
			Pattern: internal.StrPtr("test1"),
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_CollectCoverage(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"--coverage-text"}).
		Return(nil)

	err := phpUnit.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_PHPUnit_CollectCoverage_Pattern(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"--coverage-text", "tests/*"}).
		Return(nil)

	err := phpUnit.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: info,
			Pattern:     internal.StrPtr("tests/*"),
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPHPUnit_CollectCoverage_FormatUnsupported(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	data := []internal.CoverageReportFormat{
		"test",
		internal.CoverageReportFormatJson,
	}

	for _, format := range data {
		t.Run(string(format), func(t *testing.T) {
			err := phpUnit.CollectCoverage(
				t.Context(),
				internal.ProjectCoverageCollectorOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.CoverageReportFormat]{
						Format: format,
					},
					Pattern: internal.StrPtr("test1"),
				},
			)
			require.EqualError(t, err, "unsupported coverage report format: "+string(format))

			exec.AssertExpectations(t)
		})
	}
}

func TestPHPUnit_CollectCoverage_FormatSupported(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewPHPUnit(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	data := []struct {
		Format   internal.CoverageReportFormat
		Filename *string
		Args     []string
	}{
		{
			Format: internal.CoverageReportFormatRaw,
			Args:   []string{"--coverage-text"},
		},
		{
			Format:   internal.CoverageReportFormatCobertura,
			Filename: internal.StrPtr("cobertura.xml"),
			Args:     []string{"--coverage-cobertura", "cobertura.xml"},
		},
		{
			Format:   internal.CoverageReportFormatCrap4j,
			Filename: internal.StrPtr("crap4j.xml"),
			Args:     []string{"--coverage-crap4j", "crap4j.xml"},
		},
		{
			Format:   internal.CoverageReportFormatHtml,
			Filename: internal.StrPtr("coverage.html"),
			Args:     []string{"--coverage-html", "coverage.html"},
		},
		{
			Format:   internal.CoverageReportFormatXml,
			Filename: internal.StrPtr("coverage.xml"),
			Args:     []string{"--coverage-xml", "coverage.xml"},
		},
	}

	for _, values := range data {
		t.Run(string(values.Format), func(t *testing.T) {
			exec.OnRun("test", values.Args).
				Return(nil)

			err := phpUnit.CollectCoverage(
				t.Context(),
				internal.ProjectCoverageCollectorOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.CoverageReportFormat]{
						Format:   values.Format,
						Filename: values.Filename,
					},
				},
			)
			require.NoError(t, err)

			exec.AssertExpectations(t)
		})
	}
}
