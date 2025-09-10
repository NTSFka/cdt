package tool

import (
	"cdt/internal"
	"context"
)

type PHP struct {
	internal.ExecutableTool
}

// DetectPHP create a tool for php
func DetectPHP(ctx context.Context, environment internal.Environment) *PHP {
	return NewPHP(func() *internal.Executable {
		return environment.FindExecutable(ctx, "php")
	})
}

// NewPHP creates a php tool from a custom executable
func NewPHP(detect func() *internal.Executable) *PHP {
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

func (p *PHP) RunTarget(ctx context.Context, info internal.ProjectInfo, target string, args []string) error {
	return p.RunForProject(ctx, info, append([]string{"-f", target}, args...))
}
