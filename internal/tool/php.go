package tool

import (
	"cdt/internal"
)

type PHP struct {
	internal.ExecutableTool
}

// DetectPHP create a tool for php
func DetectPHP(environment internal.Environment) *PHP {
	return NewPHP(func() *internal.Executable {
		return environment.FindExecutable("php")
	})
}

// NewPHP creates a php tool from a custom executable
func NewPHP(detect func() *internal.Executable) *PHP {
	return &PHP{
		ExecutableTool: internal.MakeExecutableTool(
			"php",
			"PHP",
			" popular general-purpose scripting language that is especially suited to web development.",
			detect,
		),
	}
}

func (p *PHP) RunTarget(project internal.Project, target string, args []string) error {
	return p.RunForProject(project, append([]string{"-f", target}, args...))
}
