package tool

import (
	"context"

	"cdt/internal"
)

type PHPUnit struct {
	internal.ExecutableTool
}

// DetectPHPUnit create a tool for phpunit.
func DetectPHPUnit(ctx context.Context, environment internal.Environment) *PHPUnit {
	return NewPHPUnit(func() (*internal.Executable, error) {
		// Detect composer vendor
		if executable, err := environment.FindExecutable(ctx, "vendor/bin/phpunit"); executable != nil {
			return executable, nil
		} else if err != nil {
			return nil, err
		}

		// Detect unversioned (system default)
		if executable, err := environment.FindExecutable(ctx, "phpunit"); executable != nil {
			return executable, nil
		} else if err != nil {
			return nil, err
		}

		return nil, nil
	})
}

// NewPHPUnit creates a phpunit tool from a custom executable.
func NewPHPUnit(detect internal.ExecutableToolDetectFunc) *PHPUnit {
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

func (p *PHPUnit) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, pattern))
}
