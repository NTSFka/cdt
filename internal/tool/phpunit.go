package tool

import (
	"context"

	"cdt/internal"
)

const IdPHPUnit = "phpunit"

type PHPUnit struct {
	internal.ExecutableTool
}

// DetectPHPUnit create a tool for phpunit.
func DetectPHPUnit(
	ctx context.Context,
	config internal.ConfigTools,
	environment internal.Environment,
) *PHPUnit {
	if path, ok := config[IdPHPUnit]; ok {
		return NewPHPUnit(func() (*internal.Executable, error) {
			return environment.FindExecutable(ctx, path)
		})
	}

	return NewPHPUnit(internal.DetectExecutableChain(
		[]string{"vendor/bin/phpunit", "phpunit"},
		func(name string) (*internal.Executable, error) {
			return environment.FindExecutable(ctx, name)
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
	return p.RunForProject(ctx, options.ProjectInfo, options.ExtraArgs)
}

func (p *PHPUnit) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, pattern))
}
