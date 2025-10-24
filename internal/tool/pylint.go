package tool

import (
	"cdt/internal"
	"context"
)

type Pylint struct {
	internal.ExecutableTool
}

// DetectPylint create a tool for pylint.
func DetectPylint(ctx context.Context, environment internal.Environment) *Pylint {
	return NewPylint(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, "pylint")
	})
}

// NewPylint creates a pylint tool from a custom executable.
func NewPylint(detect internal.ExecutableToolDetectFunc) *Pylint {
	return &Pylint{
		ExecutableTool: internal.MakeExecutableTool(
			"pylint",
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
