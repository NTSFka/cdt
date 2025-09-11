package tool

import (
	"cdt/internal"
	"context"
)

type PHPUnit struct {
	internal.ExecutableTool
}

// DetectPHPUnit create a tool for phpunit
func DetectPHPUnit(ctx context.Context, environment internal.Environment) *PHPUnit {
	return NewPHPUnit(func() *internal.Executable {
		// Detect composer vendor
		if executable := environment.FindExecutable(ctx, "vendor/bin/phpunit"); executable != nil {
			return executable
		}

		// Detect unversioned (system default)
		if executable := environment.FindExecutable(ctx, "phpunit"); executable != nil {
			return executable
		}

		return nil
	})
}

// NewPHPUnit creates a phpunit tool from a custom executable
func NewPHPUnit(detect func() *internal.Executable) *PHPUnit {
	return &PHPUnit{
		ExecutableTool: internal.MakeExecutableTool(
			"phpunit",
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

func (p *PHPUnit) TestPattern(ctx context.Context, options internal.ProjectTesterOptions, pattern string) error {
	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, pattern))
}
