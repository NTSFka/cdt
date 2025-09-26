package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

type Ruff struct {
	internal.ExecutableTool
}

// DetectRuff create a tool for ruff.
func DetectRuff(ctx context.Context, environment internal.Environment) *Ruff {
	return NewRuff(func() *internal.Executable {
		return environment.FindExecutable(ctx, "ruff")
	})
}

// NewRuff creates a ruff tool from a custom executable.
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

func (r *Ruff) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return r.RunForProject(ctx, options.ProjectInfo, append([]string{"check"}, options.ExtraArgs...))
}

func (r *Ruff) LintFiles(ctx context.Context, options internal.ProjectLinterOptions, filenames []string) error {
	paths := r.buildPaths(options.Directory, filenames)

	return r.RunForProject(ctx, options.ProjectInfo, append(append([]string{"check"}, options.ExtraArgs...), paths...))
}

func (r *Ruff) FormatAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return r.RunForProject(ctx, options.ProjectInfo, append([]string{"format"}, options.ExtraArgs...))
}

func (r *Ruff) FormatFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	paths := r.buildPaths(options.Directory, filenames)

	return r.RunForProject(ctx, options.ProjectInfo, append(append([]string{"format"}, options.ExtraArgs...), paths...))
}

func (r *Ruff) FormatCheckAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return r.RunForProject(ctx, options.ProjectInfo, append([]string{"format", "--check"}, options.ExtraArgs...))
}

func (r *Ruff) FormatCheckFiles(ctx context.Context, options internal.ProjectFormatterOptions, filenames []string) error {
	paths := r.buildPaths(options.Directory, filenames)

	return r.RunForProject(ctx, options.ProjectInfo, append(append([]string{"format", "--check"}, options.ExtraArgs...), paths...))
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
