package tool

import (
	"cdt/internal"
	"path/filepath"
)

type Black struct {
	internal.ExecutableTool
}

// DetectBlack create a tool for black
func DetectBlack(environment internal.Environment) *Black {
	return NewBlack(func() *internal.Executable {
		return environment.FindExecutable("black")
	})
}

// NewBlack creates a black tool from a custom executable
func NewBlack(detect func() *internal.Executable) *Black {
	return &Black{
		ExecutableTool: internal.MakeExecutableTool(
			"black",
			"Black",
			"Black is the uncompromising Python code formatter.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagFormat},
			detect,
		),
	}
}

func (b *Black) buildPaths(directory string, filenames []string) []string {
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

func (b *Black) FormatAll(project internal.Project, args []string) error {
	return b.RunForProject(project, args)
}

func (b *Black) FormatFiles(project internal.Project, filenames []string, args []string) error {
	paths := b.buildPaths(project.Directory, filenames)

	return b.RunForProject(project, append(args, paths...))
}

func (b *Black) FormatCheckAll(project internal.Project, args []string) error {
	return b.RunForProject(project, append([]string{"--check"}, args...))
}

func (b *Black) FormatCheckFiles(project internal.Project, filenames []string, args []string) error {
	paths := b.buildPaths(project.Directory, filenames)

	return b.RunForProject(project, append(append([]string{"--check"}, args...), paths...))
}
