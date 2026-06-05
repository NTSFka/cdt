package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdFlake8 = "flake8"

type Flake8 struct {
	internal.ExecutableTool
}

// DetectFlake8 create a tool for flake8.
func DetectFlake8(
	ctx context.Context,
	options DetectOptions,
) *Flake8 {
	return NewFlake8(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdFlake8, "flake8"))
	})
}

// NewFlake8 creates a flake8 tool from a custom executable.
func NewFlake8(detect internal.ExecutableToolDetectFunc) *Flake8 {
	return &Flake8{
		ExecutableTool: internal.MakeExecutableTool(
			IdFlake8,
			"Flake8",
			"Flake8 is a wrapper around these tools: PyFlakes, pycodestyle, Ned Batchelder's McCabe script.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *Flake8) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := options.ExtraArgs

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

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
