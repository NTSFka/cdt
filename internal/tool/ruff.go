package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdRuff = "ruff"

type Ruff struct {
	internal.ExecutableTool
}

// DetectRuff create a tool for ruff.
func DetectRuff(
	ctx context.Context,
	options internal.DetectOptions,
) *Ruff {
	return NewRuff(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdRuff, "ruff"))
	})
}

// NewRuff creates a ruff tool from a custom executable.
func NewRuff(detect internal.ExecutableToolDetectFunc) *Ruff {
	return &Ruff{
		ExecutableTool: internal.MakeExecutableTool(
			IdRuff,
			"Ruff",
			"An extremely fast Python linter and code formatter, written in Rust.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint, internal.ToolTagFormat},
			detect,
		),
	}
}

func (r *Ruff) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"check"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	if options.Output.Reporter != nil {
		return fmt.Errorf("output to reporter not supported")
	}

	return r.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (r *Ruff) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	args := []string{"format"}

	if options.CheckOnly {
		args = append(args, "--check")
	}

	args = append(args, options.ExtraArgs...)

	if options.Filenames != nil {
		args = append(args, *options.Filenames...)
	}

	return r.RunForProject(ctx, options.ProjectInfo, args)
}
