package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

const IdPHPParallelLint = "php-parallel-lint"

type PHPParallelLint struct {
	internal.ExecutableTool
}

// DetectPHPParallelLint create a tool for php-parallel-lint.
func DetectPHPParallelLint(
	ctx context.Context,
	options internal.DetectOptions,
) *PHPParallelLint {
	if path, ok := options.ToolsPaths[IdPHPParallelLint]; ok {
		return NewPHPParallelLint(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewPHPParallelLint(internal.DetectExecutableChain(
		[]string{
			filepath.Join(options.ProjectDirectory, "vendor/bin/parallel-lint"),
			"parallel-lint",
		},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPParallelLint creates a php-parallel-lint tool from a custom executable.
func NewPHPParallelLint(detect internal.ExecutableToolDetectFunc) *PHPParallelLint {
	return &PHPParallelLint{
		ExecutableTool: internal.MakeExecutableTool(
			IdPHPParallelLint,
			"PHP Parallel Lint",
			"This application checks the syntax of PHP files in parallel.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *PHPParallelLint) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := appendFiles(options.ExtraArgs, options.Filenames, nil)

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
