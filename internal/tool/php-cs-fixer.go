package tool

import (
	"cdt/internal"
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

func (p *PHPCSFixer) FormatAll(project internal.Project, args []string) error {
	return p.RunForProject(project, append([]string{"fix"}, args...))
}

func (p *PHPCSFixer) FormatFiles(project internal.Project, filenames []string, args []string) error {
	paths := p.buildPaths(project.RootDirectory(), filenames)

	return p.RunForProject(project, append(append([]string{"fix"}, args...), paths...))
}

func (p *PHPCSFixer) FormatCheckAll(project internal.Project, args []string) error {
	return p.RunForProject(project, append([]string{"fix", "--dry-run"}, args...))
}

func (p *PHPCSFixer) FormatCheckFiles(project internal.Project, filenames []string, args []string) error {
	paths := p.buildPaths(project.RootDirectory(), filenames)

	return p.RunForProject(project, append(append([]string{"fix", "--dry-run"}, args...), paths...))
}
