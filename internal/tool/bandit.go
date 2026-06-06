package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdBandit = "bandit"

type Bandit struct {
	internal.ExecutableTool
}

// DetectBandit create a tool for bandit.
func DetectBandit(
	ctx context.Context,
	options DetectOptions,
) *Bandit {
	return NewBandit(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdBandit, "bandit"))
	})
}

// NewBandit creates a bandit tool from a custom executable.
func NewBandit(detect internal.ExecutableToolDetectFunc) *Bandit {
	return &Bandit{
		ExecutableTool: internal.MakeExecutableTool(
			IdBandit,
			"Bandit",
			"A security linter from PyCQA.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (b *Bandit) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := options.ExtraArgs

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	} else {
		args = append(args, "*")
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
		return b.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return b.RunForProject(ctx, options.ProjectInfo, args)
}
