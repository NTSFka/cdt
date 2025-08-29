package tool

import (
	"cdt/internal"
)

type PHPUnit struct {
	internal.ExecutableTool
}

// DetectPHPUnit create a tool for phpunit
func DetectPHPUnit(environment internal.Environment) *PHPUnit {
	return NewPHPUnit(func() *internal.Executable {
		// Detect composer vendor
		if executable := environment.FindExecutable("vendor/bin/phpunit"); executable != nil {
			return executable
		}

		// Detect unversioned (system default)
		if executable := environment.FindExecutable("phpunit"); executable != nil {
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

func (p *PHPUnit) TestAll(project internal.Project, args []string) error {
	return p.RunForProject(project, args)
}

func (p *PHPUnit) Test(project internal.Project, pattern string, args []string) error {
	return p.RunForProject(project, append(args, pattern))
}
