package tool

import (
	"context"
	"path/filepath"

	"cdt/internal"
)

type Pylint struct {
	internal.ExecutableTool
}

// DetectPylint create a tool for pylint.
func DetectPylint(ctx context.Context, environment internal.Environment) *Pylint {
	return NewPylint(func() *internal.Executable {
		return environment.FindExecutable(ctx, "pylint")
	})
}

// NewPylint creates a pylint tool from a custom executable.
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

func (p *Pylint) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return p.RunForProject(ctx, options.ProjectInfo, append([]string{"*"}, options.ExtraArgs...))
}

func (p *Pylint) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	paths := p.buildPaths(options.Directory, filenames)

	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, paths...))
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
