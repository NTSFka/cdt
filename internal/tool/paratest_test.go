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

func TestParaTest_DetectParaTest_Composer(t *testing.T) {
	paraTest := test.RunDetectToolLastFound(
		t,
		paraTestDetect,
		[]string{filepath.Join("vendor", "bin", "paratest")},
	)

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_DetectParaTest_System(t *testing.T) {
	paraTest := test.RunDetectToolLastFound(
		t,
		paraTestDetect,
		[]string{filepath.Join("vendor", "bin", "paratest"), "paratest"},
	)

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_DetectParaTest_NotFound(t *testing.T) {
	paraTest := test.RunDetectToolNotFound(
		t,
		paraTestDetect,
		[]string{filepath.Join("vendor", "bin", "paratest"), "paratest"},
	)

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_DetectParaTest_Config(t *testing.T) {
	paraTest := test.RunDetectToolConfig(t, paraTestDetect, "paratest")

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_ParaTest_RunTests(t *testing.T) {
	test.InvokeRunTestsSuccess(
		t,
		paraTestTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{},
	)
}

func TestParaTest_ParaTest_RunTests_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeRunTestsSuccess(
		t,
		paraTestTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{pattern},
	)
}

func TestParaTest_RunTests_Pattern_FormatUnsupported(t *testing.T) {
	dataSet := []internal.TestsReportFormat{
		"test",
		internal.TestsReportFormatJson,
		internal.TestsReportFormatCtrf,
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeRunTestsOutputFormatUnsupported(
				t,
				paraTestTester,
				format,
				internal.StrPtr("test1"),
			)
		})
	}
}

func TestParaTest_RunTests_Pattern_FormatSupported(t *testing.T) {
	dataSet := []struct {
		Format   internal.TestsReportFormat
		Filename *string
		Args     []string
	}{
		{
			Format:   internal.TestsReportFormatRawEvents,
			Filename: internal.StrPtr("paratest-raw-events.txt"),
			Args:     []string{"--log-events-text", "paratest-raw-events.txt"},
		},
		{
			Format:   internal.TestsReportFormatJUnit,
			Filename: internal.StrPtr("paratest-junit.xml"),
			Args:     []string{"--log-junit", "paratest-junit.xml"},
		},
		{
			Format:   internal.TestsReportFormatTeamCity,
			Filename: internal.StrPtr("paratest-teamcity.xml"),
			Args:     []string{"--log-teamcity", "paratest-teamcity.xml"},
		},
	}

	for _, values := range dataSet {
		t.Run(string(values.Format), func(t *testing.T) {
			test.InvokeRunTestsSuccess(
				t,
				paraTestTester,
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

func TestParaTest_ParaTest_CollectCoverage(t *testing.T) {
	test.InvokeCollectCoverageSuccess(
		t,
		paraTestCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{"--coverage-text"},
	)
}

func TestParaTest_ParaTest_CollectCoverage_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeCollectCoverageSuccess(
		t,
		paraTestCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{"--coverage-text", pattern},
	)
}

func TestParaTest_CollectCoverage_FormatUnsupported(t *testing.T) {
	dataSet := []internal.CoverageReportFormat{
		"test",
		internal.CoverageReportFormatJson,
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeCollectCoverageOutputFormatUnsupported(
				t,
				paraTestCoverageCollector,
				format,
				internal.StrPtr("test1"),
			)
		})
	}
}

func TestParaTest_CollectCoverage_FormatSupported(t *testing.T) {
	dataSet := []struct {
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

	for _, values := range dataSet {
		t.Run(string(values.Format), func(t *testing.T) {
			test.InvokeCollectCoverageSuccess(
				t,
				paraTestCoverageCollector,
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

func paraTestDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectParaTest(ctx, options)
}

func paraTestTester(executable func() (*internal.Executable, error)) internal.ProjectTester {
	return tool.NewParaTest(executable)
}

func paraTestCoverageCollector(
	executable func() (*internal.Executable, error),
) internal.ProjectCoverageCollector {
	return tool.NewParaTest(executable)
}
