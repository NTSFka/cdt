package tool

import (
	"cdt/internal"
	"context"
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

func (r *Ruff) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return r.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"check"}, options.ExtraArgs...),
	)
}

func (r *Ruff) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	return r.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"check"}, options.ExtraArgs...), filenames...),
	)
}

func (r *Ruff) FormatAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return r.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"format"}, options.ExtraArgs...),
	)
}

func (r *Ruff) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	return r.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"format"}, options.ExtraArgs...), filenames...),
	)
}

func (r *Ruff) FormatCheckAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return r.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"format", "--check"}, options.ExtraArgs...),
	)
}

func (r *Ruff) FormatCheckFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	return r.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"format", "--check"}, options.ExtraArgs...), filenames...),
	)
}
