package tool

import (
	"cdt/internal"
	"context"
)

const IdBlack = "black"

type Black struct {
	internal.ExecutableTool
}

// DetectBlack create a tool for black.
func DetectBlack(
	ctx context.Context,
	options internal.DetectOptions,
) *Black {
	return NewBlack(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdBlack, "black"))
	})
}

// NewBlack creates a black tool from a custom executable.
func NewBlack(detect internal.ExecutableToolDetectFunc) *Black {
	return &Black{
		ExecutableTool: internal.MakeExecutableTool(
			IdBlack,
			"Black",
			"Black is the uncompromising Python code formatter.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagFormat},
			detect,
		),
	}
}

func (b *Black) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	var args []string

	if options.CheckOnly {
		args = append(args, "--check")
	}

	args = append(args, options.ExtraArgs...)

	if options.Filenames != nil {
		args = append(args, *options.Filenames...)
	} else {
		args = append(args, ".")
	}

	return b.RunForProject(ctx, options.ProjectInfo, args)
}
