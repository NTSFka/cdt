package tool

import (
	"cdt/internal"
	"context"
)

type Black struct {
	internal.ExecutableTool
}

// DetectBlack create a tool for black.
func DetectBlack(ctx context.Context, environment internal.Environment) *Black {
	return NewBlack(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, "black")
	})
}

// NewBlack creates a black tool from a custom executable.
func NewBlack(detect internal.ExecutableToolDetectFunc) *Black {
	return &Black{
		ExecutableTool: internal.MakeExecutableTool(
			"black",
			"Black",
			"Black is the uncompromising Python code formatter.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagFormat},
			detect,
		),
	}
}

func (b *Black) FormatAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return b.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "."))
}

func (b *Black) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	return b.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, filenames...))
}

func (b *Black) FormatCheckAll(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	return b.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"--check", "."}, options.ExtraArgs...),
	)
}

func (b *Black) FormatCheckFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	return b.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"--check"}, options.ExtraArgs...), filenames...),
	)
}
