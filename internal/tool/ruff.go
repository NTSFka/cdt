package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

type Ruff struct {
	internal.ExecutableTool
}

// DetectRuff create a tool for ruff
func DetectRuff(ctx context.Context, environment internal.Environment) *Ruff {
	return NewRuff(func() *internal.Executable {
		return environment.FindExecutable(ctx, "ruff")
	})
}

// NewRuff creates a ruff tool from a custom executable
func NewRuff(detect func() *internal.Executable) *Ruff {
	return &Ruff{
		ExecutableTool: internal.MakeExecutableTool(
			"ruff",
			"Ruff",
			"An extremely fast Python linter and code formatter, written in Rust.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint, internal.ToolTagFormat},
			detect,
		),
	}
}

func (r *Ruff) buildPaths(directory string, filenames []string) []string {
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

func (r *Ruff) LintAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return r.RunForProject(ctx, info, append([]string{"check"}, args...))
}

func (r *Ruff) LintFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	paths := r.buildPaths(info.Directory, filenames)

	return r.RunForProject(ctx, info, append(append([]string{"check"}, args...), paths...))
}

func (r *Ruff) FormatAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return r.RunForProject(ctx, info, append([]string{"format"}, args...))
}

func (r *Ruff) FormatFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	paths := r.buildPaths(info.Directory, filenames)

	return r.RunForProject(ctx, info, append(append([]string{"format"}, args...), paths...))
}

func (r *Ruff) FormatCheckAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return r.RunForProject(ctx, info, append([]string{"format", "--check"}, args...))
}

func (r *Ruff) FormatCheckFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	paths := r.buildPaths(info.Directory, filenames)

	return r.RunForProject(ctx, info, append(append([]string{"format", "--check"}, args...), paths...))
}
