package tool

import (
	"context"
	"fmt"
	"path/filepath"

	"cdt/internal"
)

const IdPHPUnit = "phpunit"

type PHPUnit struct {
	internal.ExecutableTool
}

// DetectPHPUnit create a tool for phpunit.
func DetectPHPUnit(
	ctx context.Context,
	options DetectOptions,
) *PHPUnit {
	if path, ok := options.ToolsPaths[IdPHPUnit]; ok {
		return NewPHPUnit(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewPHPUnit(internal.DetectExecutableChain(
		[]string{filepath.Join(options.ProjectDirectory, "vendor/bin/phpunit"), "phpunit"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPUnit creates a phpunit tool from a custom executable.
func NewPHPUnit(detect internal.ExecutableToolDetectFunc) *PHPUnit {
	return &PHPUnit{
		ExecutableTool: internal.MakeExecutableTool(
			IdPHPUnit,
			"PHPUnit",
			"PHPUnit is a programmer-oriented testing framework for PHP.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagTest},
			detect,
		),
	}
}

func (p *PHPUnit) RunTests(ctx context.Context, options internal.ProjectTesterOptions) error {
	args := options.ExtraArgs

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	filename := internal.DefaultString(options.Output.Filename, "php://stdout")

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
	case internal.TestsReportFormatRawEvents:
		args = append(args, "--log-events-text", filename)
	case internal.TestsReportFormatJUnit:
		args = append(args, "--log-junit", filename)
	case internal.TestsReportFormatTeamCity:
		args = append(args, "--log-teamcity", filename)
	default:
		return fmt.Errorf("unsupported report format: %s", options.Output.Format)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
}

func (p *PHPUnit) CollectCoverage(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	args := options.ExtraArgs

	filename := internal.DefaultString(options.Output.Filename, "php://stdout")

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.CoverageReportFormatDefault:
		fallthrough
	case internal.CoverageReportFormatRaw:
		args = append(args, "--coverage-text")
	case internal.CoverageReportFormatCobertura:
		args = append(args, "--coverage-cobertura", filename)
	case internal.CoverageReportFormatCrap4j:
		args = append(args, "--coverage-crap4j", filename)
	case internal.CoverageReportFormatHtml:
		args = append(args, "--coverage-html", filename)
	case internal.CoverageReportFormatXml:
		args = append(args, "--coverage-xml", filename)
	default:
		return fmt.Errorf("unsupported coverage report format: %s", options.Output.Format)
	}

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	return p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		args,
	)
}
