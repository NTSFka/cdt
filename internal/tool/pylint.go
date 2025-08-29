package tool

import (
	"cdt/internal"
	"path/filepath"
)

type Pylint struct {
	internal.ExecutableTool
}

// DetectPylint create a tool for pylint
func DetectPylint(environment internal.Environment) *Pylint {
	return NewPylint(func() *internal.Executable {
		return environment.FindExecutable("pylint")
	})
}

// NewPylint creates a pylint tool from a custom executable
func NewPylint(detect func() *internal.Executable) *Pylint {
	return &Pylint{
		ExecutableTool: internal.MakeExecutableTool(
			"pylint",
			"Pylint",
			"Pylint is a static code analyser.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *Pylint) buildPaths(directory string, filenames []string) []string {
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

func (p *Pylint) LintAll(project internal.Project, args []string) error {
	return p.RunForProject(project, append([]string{"*"}, args...))
}

func (p *Pylint) LintFiles(project internal.Project, filenames []string, args []string) error {
	paths := p.buildPaths(project.RootDirectory(), filenames)

	return p.RunForProject(project, append(args, paths...))
}
