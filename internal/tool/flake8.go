package tool

import (
	"cdt/internal"
	"path/filepath"
)

type Flake8 struct {
	internal.ExecutableTool
}

// DetectFlake8 create a tool for flake8
func DetectFlake8(environment internal.Environment) *Flake8 {
	return NewFlake8(func() *internal.Executable {
		return environment.FindExecutable("flake8")
	})
}

// NewFlake8 creates a flake8 tool from a custom executable
func NewFlake8(detect func() *internal.Executable) *Flake8 {
	return &Flake8{
		ExecutableTool: internal.MakeExecutableTool(
			"flake8",
			"Flake8",
			"Flake8 is a wrapper around these tools: PyFlakes, pycodestyle, Ned Batchelder's McCabe script.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *Flake8) buildPaths(directory string, filenames []string) []string {
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

func (p *Flake8) LintAll(info internal.ProjectInfo, args []string) error {
	return p.RunForProject(info, args)
}

func (p *Flake8) LintFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	paths := p.buildPaths(info.Directory, filenames)

	return p.RunForProject(info, append(args, paths...))
}
