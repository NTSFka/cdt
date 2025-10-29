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

func (p *Pylint) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return p.RunForProject(ctx, options.ProjectInfo, append([]string{"*"}, options.ExtraArgs...))
}

func (p *Pylint) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, filenames...))
}
