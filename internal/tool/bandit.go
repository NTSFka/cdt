package tool

import (
	"context"
	"path/filepath"

	"cdt/internal"
)

type Bandit struct {
	internal.ExecutableTool
}

// DetectBandit create a tool for bandit.
func DetectBandit(ctx context.Context, environment internal.Environment) *Bandit {
	return NewBandit(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, "bandit")
	})
}

// NewBandit creates a bandit tool from a custom executable.
func NewBandit(detect internal.ExecutableToolDetectFunc) *Bandit {
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

func (b *Bandit) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return b.RunForProject(ctx, options.ProjectInfo, append([]string{"*"}, options.ExtraArgs...))
}

func (b *Bandit) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	paths := b.buildPaths(options.Directory, filenames)

	return b.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, paths...))
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
