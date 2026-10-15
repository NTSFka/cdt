package tool_test

import (
	"context"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestPHPUnit_DetectPHPUnit_Composer(t *testing.T) {
	phpUnit := test.RunDetectToolLastFound(
		t,
		phpUnitDetect,
		[]string{filepath.Join("vendor", "bin", "phpunit")},
	)

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_DetectPHPUnit_System(t *testing.T) {
	phpUnit := test.RunDetectToolLastFound(
		t,
		phpUnitDetect,
		[]string{filepath.Join("vendor", "bin", "phpunit"), "phpunit"},
	)

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_DetectPHPUnit_NotFound(t *testing.T) {
	phpUnit := test.RunDetectToolNotFound(
		t,
		phpUnitDetect,
		[]string{filepath.Join("vendor", "bin", "phpunit"), "phpunit"},
	)

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_DetectPHPUnit_Config(t *testing.T) {
	phpUnit := test.RunDetectToolConfig(t, phpUnitDetect, "phpunit")

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_PHPUnit_RunTests(t *testing.T) {
	test.InvokeRunTestsSuccess(
		t,
		phpUnitTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{},
	)
}

func TestPHPUnit_PHPUnit_RunTests_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeRunTestsSuccess(
		t,
		phpUnitTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{pattern},
	)
}

func TestPHPUnit_RunTests_Pattern_FormatUnsupported(t *testing.T) {
	dataSet := []internal.TestsReportFormat{
		"test",
		internal.TestsReportFormatJson,
		internal.TestsReportFormatCtrf,
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeRunTestsOutputFormatUnsupported(
				t,
				phpUnitTester,
				format,
				internal.StrPtr("test1"),
			)
		})
	}
}

func TestPHPUnit_RunTests_Pattern_FormatSupported(t *testing.T) {
	dataSet := []struct {
		Format   internal.TestsReportFormat
		Filename *string
		Args     []string
	}{
		{
			Format:   internal.TestsReportFormatRawEvents,
			Filename: internal.StrPtr("phpunit-raw-events.txt"),
			Args:     []string{"--log-events-text", "phpunit-raw-events.txt"},
		},
		{
			Format:   internal.TestsReportFormatJUnit,
			Filename: internal.StrPtr("phpunit-junit.xml"),
			Args:     []string{"--log-junit", "phpunit-junit.xml"},
		},
		{
			Format:   internal.TestsReportFormatTeamCity,
			Filename: internal.StrPtr("phpunit-teamcity.xml"),
			Args:     []string{"--log-teamcity", "phpunit-teamcity.xml"},
		},
	}

	for _, values := range dataSet {
		t.Run(string(values.Format), func(t *testing.T) {
			test.InvokeRunTestsSuccess(
				t,
				phpUnitTester,
				internal.ProjectTesterOptions{
					ProjectInfo: internal.ProjectInfo{Directory: "."},
					Output: internal.OutputOptions[internal.TestsReportFormat]{
						Format:   values.Format,
						Filename: values.Filename,
					},
				},
				values.Args,
			)
		})
	}
}

func TestPHPUnit_PHPUnit_CollectCoverage(t *testing.T) {
	test.InvokeCollectCoverageSuccess(
		t,
		phpUnitCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{"--coverage-text"},
	)
}

func TestPHPUnit_PHPUnit_CollectCoverage_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeCollectCoverageSuccess(
		t,
		phpUnitCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{"--coverage-text", pattern},
	)
}

func TestPHPUnit_CollectCoverage_FormatUnsupported(t *testing.T) {
	dataSet := []internal.CoverageReportFormat{
		"test",
		internal.CoverageReportFormatJson,
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeCollectCoverageOutputFormatUnsupported(
				t,
				phpUnitCoverageCollector,
				format,
				internal.StrPtr("test1"),
			)
		})
	}
}

func TestPHPUnit_CollectCoverage_FormatSupported(t *testing.T) {
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
			test.InvokeCollectCoverageSuccess(
				t,
				phpUnitCoverageCollector,
				internal.ProjectCoverageCollectorOptions{
					ProjectInfo: internal.ProjectInfo{Directory: "."},
					Output: internal.OutputOptions[internal.CoverageReportFormat]{
						Format:   values.Format,
						Filename: values.Filename,
					},
				},
				values.Args,
			)
		})
	}
}

func phpUnitDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHPUnit(ctx, options)
}

func phpUnitTester(executable func() (*internal.Executable, error)) internal.ProjectTester {
	return tool.NewPHPUnit(executable)
}

func phpUnitCoverageCollector(
	executable func() (*internal.Executable, error),
) internal.ProjectCoverageCollector {
	return tool.NewPHPUnit(executable)
}
