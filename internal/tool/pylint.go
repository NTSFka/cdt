package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdPylint = "pylint"

type Pylint struct {
	internal.ExecutableTool
}

// DetectPylint create a tool for pylint.
func DetectPylint(
	ctx context.Context,
	options internal.DetectOptions,
) *Pylint {
	return NewPylint(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdPylint, "pylint"))
	})
}

// NewPylint creates a pylint tool from a custom executable.
func NewPylint(detect internal.ExecutableToolDetectFunc) *Pylint {
	return &Pylint{
		ExecutableTool: internal.MakeExecutableTool(
			IdPylint,
			"Pylint",
			"Pylint is a static code analyser.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *Pylint) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := appendFiles(options.ExtraArgs, options.Filenames, new("*"))

	if options.Output.Reporter != nil {
		return fmt.Errorf("output to reporter not supported")
	}

	return p.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}
