package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

const IdMago = "mago"

type Mago struct {
	internal.ExecutableTool
}

// DetectMago create a tool for mago.
func DetectMago(
	ctx context.Context,
	options internal.DetectOptions,
) *Mago {
	if path, ok := options.ToolsPaths[IdMago]; ok {
		return NewMago(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewMago(internal.DetectExecutableChain(
		[]string{filepath.Join(options.ProjectDirectory, "vendor/bin/mago"), "mago"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewMago creates a mago tool from a custom executable.
func NewMago(detect internal.ExecutableToolDetectFunc) *Mago {
	return &Mago{
		ExecutableTool: internal.MakeExecutableTool(
			IdMago,
			"Mago",
			"Mago is a static analyzer, linter, and formatter for PHP, written in Rust.",
			internal.Tags{
				internal.ToolTagPhp,
				internal.ToolTagLint,
				internal.ToolTagAnalyze,
				internal.ToolTagFormat,
			},
			detect,
		),
	}
}

func (p *Mago) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"lint"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	return p.RunForProject(ctx, options.ProjectInfo, args)
}

func (p *Mago) AnalyzeFiles(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
) error {
	args := append([]string{"analyze"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	return p.RunForProject(ctx, options.ProjectInfo, args)
}

func (p *Mago) FormatFiles(ctx context.Context, options internal.ProjectFormatterOptions) error {
	args := append([]string{"format"}, options.ExtraArgs...)

	if options.CheckOnly {
		args = append(args, "--dry-run")
	}

	args = appendFiles(args, options.Filenames, nil)

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
