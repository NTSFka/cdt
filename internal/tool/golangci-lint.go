package tool

import (
	"cdt/internal"
)

// A GolangCILint is a tool that wraps golang main tool `golangci-lint`.
type GolangCILint struct {
	internal.ExecutableTool
}

// NewGolangCILint creates a go tool from a custom executable
func NewGolangCILint(detect func() *internal.Executable) *GolangCILint {
	return &GolangCILint{
		internal.MakeExecutableTool(
			"golangci-lint",
			"Golangci-lint",
			"Smart, fast linters runner.",
			detect,
		),
	}
}

// DetectGolangCILint create golangci-lint tool can be used in the project
func DetectGolangCILint(environment internal.Environment) *GolangCILint {
	return NewGolangCILint(func() *internal.Executable {
		return environment.FindExecutable("golangci-lint")
	})
}

func (c *GolangCILint) LintAll(project internal.Project, args []string) error {
	return c.RunForProject(project, append(args, "run"))
}

func (c *GolangCILint) LintFiles(project internal.Project, filenames []string, args []string) error {
	return c.RunForProject(project, append(append(args, "run"), filenames...))
}
