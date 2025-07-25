package tool

import (
	"cdt/internal"
	"path/filepath"
)

type PHPStan struct {
	internal.ExecutableTool
}

// DetectPHPStan create a tool for phpstan
func DetectPHPStan(environment internal.Environment) *PHPStan {
	return NewPHPStan(func() *internal.Executable {
		// Detect composer vendor
		if executable := environment.FindExecutable("vendor/bin/phpstan"); executable != nil {
			return executable
		}

		// Detect unversioned (system default)
		if executable := environment.FindExecutable("phpstan"); executable != nil {
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

func (p *PHPStan) LintAll(project internal.Project, args []string) error {
	return p.RunForProject(project, append([]string{"analyse"}, args...))
}

func (p *PHPStan) LintFiles(project internal.Project, filenames []string, args []string) error {
	paths := p.buildPaths(project.RootDirectory(), filenames)

	return p.RunForProject(project, append(append([]string{"analyse"}, args...), paths...))
}
