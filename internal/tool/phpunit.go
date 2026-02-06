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
	return p.runTests(ctx, options)
}

func (p *PHPUnit) CollectCoverage(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	args := options.ExtraArgs

	args = append(args, "--coverage-text")

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

func (p *PHPUnit) runTests(
	ctx context.Context,
	options internal.ProjectTesterOptions,
) error {
	args := options.ExtraArgs

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	filename := internal.DefaultString(options.Output.Filename, "php://stdout")

	// TODO: other report formats
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
		return p.RunForProject(ctx, options.ProjectInfo, args)
	case internal.TestsReportFormatRawEvents:
		return p.RunForProject(
			ctx,
			options.ProjectInfo,
			append(args, "--log-events-text", filename),
		)
	case internal.TestsReportFormatJson:
		break
	case internal.TestsReportFormatCtrf:
		break
	case internal.TestsReportFormatJUnit:
		return p.RunForProject(
			ctx,
			options.ProjectInfo,
			append(args, "--log-junit", filename),
		)
	case internal.TestsReportFormatTeamCity:
		return p.RunForProject(
			ctx,
			options.ProjectInfo,
			append(args, "--log-teamcity", filename),
		)
	}

	return fmt.Errorf("unsupported report format: %s", options.Output.Format)
}
