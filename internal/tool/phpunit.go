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

func (p *PHPUnit) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return p.runTests(ctx, options, nil)
}

func (p *PHPUnit) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return p.runTests(ctx, options, &pattern)
}

func (p *PHPUnit) CollectCoverageAll(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	return p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		append(options.ExtraArgs, "--coverage-text"),
	)
}

func (p *PHPUnit) CollectCoveragePattern(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
	pattern string,
) error {
	return p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		append(options.ExtraArgs, "--coverage-text", pattern),
	)
}

func (p *PHPUnit) runTests(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern *string,
) error {
	// TODO: other report formats
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
		args := options.ExtraArgs

		if pattern != nil {
			args = append(args, *pattern)
		}

		return p.RunForProject(ctx, options.ProjectInfo, args)
	case internal.TestsReportFormatEvents:
		break
	case internal.TestsReportFormatJson:
		break
	case internal.TestsReportFormatCtrf:
		break
	}

	return fmt.Errorf("unsupported report format: %s", options.Output.Format)
}
