package tool

import (
	"cdt/internal"
	"path/filepath"
)

type Bandit struct {
	internal.ExecutableTool
}

// DetectBandit create a tool for bandit
func DetectBandit(environment internal.Environment) *Bandit {
	return NewBandit(func() *internal.Executable {
		return environment.FindExecutable("bandit")
	})
}

// NewBandit creates a bandit tool from a custom executable
func NewBandit(detect func() *internal.Executable) *Bandit {
	return &Bandit{
		ExecutableTool: internal.MakeExecutableTool(
			"bandit",
			"Bandit",
			"A security linter from PyCQA.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (b *Bandit) buildPaths(directory string, filenames []string) []string {
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

func (b *Bandit) LintAll(project internal.Project, args []string) error {
	return b.RunForProject(project, append([]string{"*"}, args...))
}

func (b *Bandit) LintFiles(project internal.Project, filenames []string, args []string) error {
	paths := b.buildPaths(project.Directory, filenames)

	return b.RunForProject(project, append(args, paths...))
}
