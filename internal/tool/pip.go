package tool

import (
	"cdt/internal"
	"context"
	"errors"
)

type Pip struct {
	internal.ExecutableTool
}

// DetectPip create a tool for pip
func DetectPip(environment internal.Environment) *Pip {
	return NewPip(func() *internal.Executable {
		if executable := environment.FindExecutable("pip"); executable != nil {
			return executable
		}

		if executable := environment.FindExecutable("pip3"); executable != nil {
			return executable
		}

		return nil
	})
}

// NewPip creates a pip tool from a custom executable
func NewPip(detect func() *internal.Executable) *Pip {
	return &Pip{
		ExecutableTool: internal.MakeExecutableTool(
			"pip",
			"pip",
			"pip is the package installer for Python. ",
			internal.Tags{internal.ToolTagPython, internal.ToolTagDependency},
			detect,
		),
	}
}

func (p *Pip) AddDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string, _ bool) error {
	return p.RunForProject(ctx, info, append([]string{"install"}, dependencies...))
}

func (p *Pip) RemoveDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string, _ bool) error {
	return p.RunForProject(ctx, info, append([]string{"uninstall"}, dependencies...))
}

func (p *Pip) UpdateDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string) error {
	return p.RunForProject(ctx, info, append([]string{"install", "--upgrade"}, dependencies...))
}

func (p *Pip) FetchDependencies(ctx context.Context, info internal.ProjectInfo, _ bool) error {
	return p.RunForProject(ctx, info, []string{"install", "-r", "requirements.txt"})
}

func (p *Pip) ListDependencies(ctx context.Context, info internal.ProjectInfo) error {
	return p.RunForProject(ctx, info, []string{"list"})
}

func (p *Pip) AuditDependencies(_ context.Context, _ internal.ProjectInfo) error {
	return errors.New("not supported")
}
