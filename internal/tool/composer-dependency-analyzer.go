package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

const IdComposerDependencyAnalyzer = "composer-dependency-analyzer"

type ComposerDependencyAnalyzer struct {
	internal.ExecutableTool
}

// DetectComposerDependencyAnalyzer create a tool for composer-dependency-analyzer.
func DetectComposerDependencyAnalyzer(
	ctx context.Context,
	options internal.DetectOptions,
) *ComposerDependencyAnalyzer {
	if path, ok := options.ToolsPaths[IdComposerDependencyAnalyzer]; ok {
		return NewComposerDependencyAnalyzer(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewComposerDependencyAnalyzer(internal.DetectExecutableChain(
		[]string{
			filepath.Join(options.ProjectDirectory, "vendor/bin/composer-dependency-analyzer"),
			"composer-dependency-analyzer",
		},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewComposerDependencyAnalyzer creates a composer-dependency-analyzer tool from a custom executable.
func NewComposerDependencyAnalyzer(
	detect internal.ExecutableToolDetectFunc,
) *ComposerDependencyAnalyzer {
	return &ComposerDependencyAnalyzer{
		ExecutableTool: internal.MakeExecutableTool(
			IdComposerDependencyAnalyzer,
			"ComposerDependencyAnalyzer",
			"Detects unused, shadow and misplaced composer dependencies.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *ComposerDependencyAnalyzer) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := options.ExtraArgs
	args = appendFiles(args, options.Filenames, nil)

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
