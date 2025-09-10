package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

type PHPCSFixer struct {
	internal.ExecutableTool
}

// DetectPHPCSFixer create a tool for php-cs-fixer
func DetectPHPCSFixer(environment internal.Environment) *PHPCSFixer {
	return NewPHPCSFixer(func() *internal.Executable {
		// Detect composer vendor
		if executable := environment.FindExecutable("vendor/bin/php-cs-fixer"); executable != nil {
			return executable
		}

		// Detect unversioned (system default)
		if executable := environment.FindExecutable("php-cs-fixer"); executable != nil {
			return executable
		}

		return nil
	})
}

// NewPHPCSFixer creates a php-cs-fixer tool from a custom executable
func NewPHPCSFixer(detect func() *internal.Executable) *PHPCSFixer {
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

func (p *PHPCSFixer) FormatAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return p.RunForProject(ctx, info, append([]string{"fix"}, args...))
}

func (p *PHPCSFixer) FormatFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	paths := p.buildPaths(info.Directory, filenames)

	return p.RunForProject(ctx, info, append(append([]string{"fix"}, args...), paths...))
}

func (p *PHPCSFixer) FormatCheckAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return p.RunForProject(ctx, info, append([]string{"fix", "--dry-run"}, args...))
}

func (p *PHPCSFixer) FormatCheckFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	paths := p.buildPaths(info.Directory, filenames)

	return p.RunForProject(ctx, info, append(append([]string{"fix", "--dry-run"}, args...), paths...))
}
