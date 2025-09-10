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

func (b *Black) FormatAll(info internal.ProjectInfo, args []string) error {
	return b.RunForProject(info, args)
}

func (b *Black) FormatFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	paths := b.buildPaths(info.Directory, filenames)

	return b.RunForProject(info, append(args, paths...))
}

func (b *Black) FormatCheckAll(info internal.ProjectInfo, args []string) error {
	return b.RunForProject(info, append([]string{"--check"}, args...))
}

func (b *Black) FormatCheckFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	paths := b.buildPaths(info.Directory, filenames)

	return b.RunForProject(info, append(append([]string{"--check"}, args...), paths...))
}
