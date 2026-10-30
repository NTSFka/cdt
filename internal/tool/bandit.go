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
	options internal.DetectOptions,
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
	args := appendFiles(options.ExtraArgs, options.Filenames, new("*"))

	if options.Output.Reporter != nil {
		return fmt.Errorf("output to reporter not supported")
	}

	return b.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}
