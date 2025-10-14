package tool

import (
	"context"
	"path/filepath"

	"cdt/internal"
)

type PHPStan struct {
	internal.ExecutableTool
}

// DetectPHPStan create a tool for phpstan.
func DetectPHPStan(ctx context.Context, environment internal.Environment) *PHPStan {
	return NewPHPStan(internal.DetectExecutableChain(
		[]string{"vendor/bin/phpstan", "phpstan"},
		func(name string) (*internal.Executable, error) {
			return environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPStan creates a phpstan tool from a custom executable.
func NewPHPStan(detect internal.ExecutableToolDetectFunc) *PHPStan {
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

func (p *PHPStan) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"analyse"}, options.ExtraArgs...),
	)
}

func (p *PHPStan) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	paths := p.buildPaths(options.Directory, filenames)

	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"analyse"}, options.ExtraArgs...), paths...),
	)
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
