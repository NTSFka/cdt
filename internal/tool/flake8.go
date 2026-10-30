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
	options internal.DetectOptions,
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

func (f *Flake8) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := appendFiles(options.ExtraArgs, options.Filenames, nil)

	if options.Output.Reporter != nil {
		return fmt.Errorf("output to reporter not supported")
	}

	return f.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}
