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
	args := appendFiles(options.ExtraArgs, options.Filenames, internal.StrPtr("*"))

	if a, err := b.argsBuildLintOutputFormat(options.Output.Format); err == nil {
		args = append(args, a...)
	} else {
		return err
	}

	if options.Output.Filename != nil {
		return b.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return b.RunForProject(ctx, options.ProjectInfo, args)
}

func (b *Bandit) argsBuildLintOutputFormat(format internal.LintReportFormat) ([]string, error) {
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
