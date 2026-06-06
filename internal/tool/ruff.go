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
	options DetectOptions,
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

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	}

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.LintReportFormatDefault:
		fallthrough
	case internal.LintReportFormatRaw:
		break
	default:
		return fmt.Errorf("unsupported report format: %s", options.Output.Format)
	}

	if options.Output.Filename != nil {
		return r.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return r.RunForProject(ctx, options.ProjectInfo, args)
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
