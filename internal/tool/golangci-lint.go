package tool

import (
	"context"
	"fmt"

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

func (g *GolangCILint) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"run"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	if a, err := g.argsBuildLintOutputFormat(options.Output.Format); err == nil {
		args = append(args, a...)
	} else {
		return err
	}

	if options.Output.Filename != nil {
		return g.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return g.RunForProject(ctx, options.ProjectInfo, args)
}

func (g *GolangCILint) FormatFiles(
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

	return g.RunForProject(ctx, options.ProjectInfo, args)
}

func (g *GolangCILint) argsBuildLintOutputFormat(
	format internal.LintReportFormat,
) ([]string, error) {
	var args []string

	// nolint: exhaustive
	switch format {
	case internal.LintReportFormatDefault:
		fallthrough
	case internal.LintReportFormatRaw:
		break
	default:
		return nil, fmt.Errorf("unsupported report format: %s", format)
	}

	return args, nil
}
