package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

const IdPHPCSFixer = "php-cs-fixer"

type PHPCSFixer struct {
	internal.ExecutableTool
}

// DetectPHPCSFixer create a tool for php-cs-fixer.
func DetectPHPCSFixer(
	ctx context.Context,
	options internal.DetectOptions,
) *PHPCSFixer {
	if path, ok := options.ToolsPaths[IdPHPCSFixer]; ok {
		return NewPHPCSFixer(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewPHPCSFixer(internal.DetectExecutableChain(
		[]string{
			filepath.Join(options.ProjectDirectory, "vendor/bin/php-cs-fixer"),
			"php-cs-fixer",
		},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPCSFixer creates a php-cs-fixer tool from a custom executable.
func NewPHPCSFixer(detect internal.ExecutableToolDetectFunc) *PHPCSFixer {
	return &PHPCSFixer{
		ExecutableTool: internal.MakeExecutableTool(
			IdPHPCSFixer,
			"PHP-CS-Fixer",
			"PHP Coding Standards Fixer",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagFormat},
			detect,
		),
	}
}

func (p *PHPCSFixer) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	args := append([]string{"fix"}, options.ExtraArgs...)

	if options.CheckOnly {
		args = append(args, "--dry-run")
	}

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
