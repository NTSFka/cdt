package tool

import (
	"cdt/internal"
	"context"
)

const IdPylint = "pylint"

type Pylint struct {
	internal.ExecutableTool
}

// DetectPylint create a tool for pylint.
func DetectPylint(
	ctx context.Context,
	options DetectOptions,
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
	args := options.ExtraArgs

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	} else {
		args = append(args, "*")
	}

	if options.Output.Filename != nil {
		return p.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
