package tool

import (
	"context"

	"cdt/internal"
)

type PHP struct {
	internal.ExecutableTool
}

// DetectPHP create a tool for php.
func DetectPHP(ctx context.Context, environment internal.Environment) *PHP {
	return NewPHP(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, "php")
	})
}

// NewPHP creates a php tool from a custom executable.
func NewPHP(detect internal.ExecutableToolDetectFunc) *PHP {
	return &PHP{
		ExecutableTool: internal.MakeExecutableTool(
			"php",
			"PHP",
			"Popular general-purpose scripting language that is especially suited to web development.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagRun},
			detect,
		),
	}
}

func (p *PHP) RunTarget(
	ctx context.Context,
	options internal.ProjectRunnerOptions,
	target string,
) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"-f", target}, options.ExtraArgs...),
	)
}
