package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdPythonCoverage = "python-coverage"

type PythonCoverage struct {
	internal.ExecutableTool
	detectRunModule func() string
}

// DetectPythonCoverage create a tool for python coverage.
func DetectPythonCoverage(
	ctx context.Context,
	options internal.DetectOptions,
) *PythonCoverage {
	return NewPythonCoverage(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(
			ctx,
			options.GetToolPath(IdPythonCoverage, "coverage"),
		)
	}, func() string {
		executable, _ := options.Environment.FindExecutable(ctx, "pytest")

		if executable != nil {
			return "pytest"
		}

		return "unittest"
	})
}

// NewPythonCoverage creates a python coverage tool from a custom executable.
func NewPythonCoverage(
	detect internal.ExecutableToolDetectFunc,
	detectRunModule func() string,
) *PythonCoverage {
	if detectRunModule == nil {
		detectRunModule = func() string { return "unittest" }
	}

	return &PythonCoverage{
		ExecutableTool: internal.MakeExecutableTool(
			IdPythonCoverage,
			"Python coverage",
			"Coverage.py is a tool for measuring code coverage of Python programs.",
			internal.Tags{internal.ToolTagPython},
			detect,
		),
		detectRunModule: detectRunModule,
	}
}

// DetectRunModule returns the module to use for test run, e.g., unittest or pytest.
func (p *PythonCoverage) DetectRunModule() string {
	return p.detectRunModule()
}

func (p *PythonCoverage) CollectCoverage(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	args := []string{"run", "-m", p.detectRunModule()}

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	filename := internal.DefaultString(options.Output.Filename, "php://stdout")

	var reportArgs []string

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.CoverageReportFormatDefault:
		fallthrough
	case internal.CoverageReportFormatRaw:
		reportArgs = append(reportArgs, "report")
	case internal.CoverageReportFormatHtml:
		reportArgs = append(reportArgs, "html", filename)
	case internal.CoverageReportFormatXml:
		reportArgs = append(reportArgs, "xml", filename)
	case internal.CoverageReportFormatJson:
		reportArgs = append(reportArgs, "json", filename)
	case internal.CoverageReportFormatLcov:
		reportArgs = append(reportArgs, "lcov", filename)
	default:
		return fmt.Errorf("unsupported coverage report format: %s", options.Output.Format)
	}

	reportArgs = append(reportArgs, options.ExtraArgs...)

	// Collect coverage
	if err := p.RunForProject(ctx, options.ProjectInfo, args); err != nil {
		return err
	}

	return p.RunForProject(ctx, options.ProjectInfo, reportArgs)
}
