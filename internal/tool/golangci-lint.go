package tool

import (
	"context"

	"cdt/internal"
)

const IdGolangCILint = "golangci-lint"

// A GolangCILint is a tool that wraps golang main tool `golangci-lint`.
type GolangCILint struct {
	internal.ExecutableTool
}

// NewGolangCILint creates a go tool from a custom executable.
func NewGolangCILint(detect internal.ExecutableToolDetectFunc) *GolangCILint {
	return &GolangCILint{
		internal.MakeExecutableTool(
			IdGolangCILint,
			"Golangci-lint",
			"Smart, fast linters runner.",
			internal.Tags{internal.ToolTagGo, internal.ToolTagLint},
			detect,
		),
	}
}

// DetectGolangCILint create golangci-lint tool can be used in the project.
func DetectGolangCILint(
	ctx context.Context,
	options DetectOptions,
) *GolangCILint {
	return NewGolangCILint(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(
			ctx,
			options.GetToolPath(IdGolangCILint, "golangci-lint"),
		)
	})
}

func (l *GolangCILint) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return l.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "run"))
}

func (l *GolangCILint) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	return l.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append(options.ExtraArgs, "run"), filenames...),
	)
}

func (l *GolangCILint) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	args := append([]string{"fmt"}, options.ExtraArgs...)

	if options.CheckOnly {
		args = append(args, "--diff")
	}

	if options.Filenames != nil {
		args = append(args, *options.Filenames...)
	}

	return l.RunForProject(ctx, options.ProjectInfo, args)
}
