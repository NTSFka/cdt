package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

type PHPStan struct {
	internal.ExecutableTool
}

// DetectPHPStan create a tool for phpstan
func DetectPHPStan(ctx context.Context, environment internal.Environment) *PHPStan {
	return NewPHPStan(func() *internal.Executable {
		// Detect composer vendor
		if executable := environment.FindExecutable(ctx, "vendor/bin/phpstan"); executable != nil {
			return executable
		}

		// Detect unversioned (system default)
		if executable := environment.FindExecutable(ctx, "phpstan"); executable != nil {
			return executable
		}

		return nil
	})
}

// NewPHPStan creates a phpstan tool from a custom executable
func NewPHPStan(detect func() *internal.Executable) *PHPStan {
	return &PHPStan{
		ExecutableTool: internal.MakeExecutableTool(
			"phpstan",
			"PHPStan",
			"Analyses source code.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *PHPStan) buildPaths(directory string, filenames []string) []string {
	var paths []string

	for _, filename := range filenames {
		if filepath.IsAbs(filename) {
			paths = append(paths, filename)
		} else {
			paths = append(paths, filepath.Join(directory, filename))
		}
	}

	return paths
}

func (p *PHPStan) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return p.RunForProject(ctx, options.ProjectInfo, append([]string{"analyse"}, options.ExtraArgs...))
}

func (p *PHPStan) LintFiles(ctx context.Context, options internal.ProjectLinterOptions, filenames []string) error {
	paths := p.buildPaths(options.Directory, filenames)

	return p.RunForProject(ctx, options.ProjectInfo, append(append([]string{"analyse"}, options.ExtraArgs...), paths...))
}
