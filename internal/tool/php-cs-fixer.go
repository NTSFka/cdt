package tool

import (
	"context"
	"path/filepath"

	"cdt/internal"
)

type PHPCSFixer struct {
	internal.ExecutableTool
}

// DetectPHPCSFixer create a tool for php-cs-fixer.
func DetectPHPCSFixer(ctx context.Context, environment internal.Environment) *PHPCSFixer {
	return NewPHPCSFixer(internal.DetectExecutableChain(
		[]string{"vendor/bin/php-cs-fixer", "php-cs-fixer"},
		func(name string) (*internal.Executable, error) {
			return environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPCSFixer creates a php-cs-fixer tool from a custom executable.
func NewPHPCSFixer(detect internal.ExecutableToolDetectFunc) *PHPCSFixer {
	return &PHPCSFixer{
		ExecutableTool: internal.MakeExecutableTool(
			"php-cs-fixer",
			"PHP-CS-Fixer",
			"PHP Coding Standards Fixer",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagFormat},
			detect,
		),
	}
}

func (p *PHPCSFixer) FormatAll(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	return p.RunForProject(ctx, options.ProjectInfo, append([]string{"fix"}, options.ExtraArgs...))
}

func (p *PHPCSFixer) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	paths := p.buildPaths(options.Directory, filenames)

	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"fix"}, options.ExtraArgs...), paths...),
	)
}

func (p *PHPCSFixer) FormatCheckAll(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"fix", "--dry-run"}, options.ExtraArgs...),
	)
}

func (p *PHPCSFixer) FormatCheckFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	paths := p.buildPaths(options.Directory, filenames)

	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"fix", "--dry-run"}, options.ExtraArgs...), paths...),
	)
}

func (p *PHPCSFixer) buildPaths(directory string, filenames []string) []string {
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
