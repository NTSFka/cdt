package tool

import (
	"context"
	"path/filepath"

	"cdt/internal"
)

type Flake8 struct {
	internal.ExecutableTool
}

// DetectFlake8 create a tool for flake8.
func DetectFlake8(ctx context.Context, environment internal.Environment) *Flake8 {
	return NewFlake8(func() *internal.Executable {
		return environment.FindExecutable(ctx, "flake8")
	})
}

// NewFlake8 creates a flake8 tool from a custom executable.
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

func (p *Flake8) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return p.RunForProject(ctx, options.ProjectInfo, options.ExtraArgs)
}

func (p *Flake8) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	paths := p.buildPaths(options.Directory, filenames)

	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, paths...))
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
