package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

type Black struct {
	internal.ExecutableTool
}

// DetectBlack create a tool for black
func DetectBlack(ctx context.Context, environment internal.Environment) *Black {
	return NewBlack(func() *internal.Executable {
		return environment.FindExecutable(ctx, "black")
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

func (b *Black) FormatAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return b.RunForProject(ctx, options.ProjectInfo, options.ExtraArgs)
}

func (b *Black) FormatFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	paths := b.buildPaths(options.Directory, filenames)

	return b.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, paths...))
}

func (b *Black) FormatCheckAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return b.RunForProject(ctx, options.ProjectInfo, append([]string{"--check"}, options.ExtraArgs...))
}

func (b *Black) FormatCheckFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	paths := b.buildPaths(options.Directory, filenames)

	return b.RunForProject(ctx, options.ProjectInfo, append(append([]string{"--check"}, options.ExtraArgs...), paths...))
}
