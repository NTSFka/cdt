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

func (p *Pylint) LintAll(info internal.ProjectInfo, args []string) error {
	return p.RunForProject(info, append([]string{"*"}, args...))
}

func (p *Pylint) LintFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	paths := p.buildPaths(info.Directory, filenames)

	return p.RunForProject(info, append(args, paths...))
}
