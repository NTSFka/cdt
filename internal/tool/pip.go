package tool

import (
	"context"
	"errors"

	"cdt/internal"
)

type Pip struct {
	internal.ExecutableTool
}

// DetectPip create a tool for pip.
func DetectPip(ctx context.Context, environment internal.Environment) *Pip {
	return NewPip(func() (*internal.Executable, error) {
		if executable, err := environment.FindExecutable(ctx, "pip"); executable != nil {
			return executable, nil
		} else if err != nil {
			return nil, err
		}

		if executable, err := environment.FindExecutable(ctx, "pip3"); executable != nil {
			return executable, nil
		} else if err != nil {
			return nil, err
		}

		return nil, nil
	})
}

// NewPip creates a pip tool from a custom executable.
func NewPip(detect internal.ExecutableToolDetectFunc) *Pip {
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

func (p *Pip) AddDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	_ bool,
) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"install"}, options.ExtraArgs...), dependencies...),
	)
}

func (p *Pip) RemoveDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	_ bool,
) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"uninstall"}, options.ExtraArgs...), dependencies...),
	)
}

func (p *Pip) UpdateDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append([]string{"install", "--upgrade"}, options.ExtraArgs...), dependencies...),
	)
}

func (p *Pip) FetchDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	_ bool,
) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"install", "-r", "requirements.txt"}, options.ExtraArgs...),
	)
}

func (p *Pip) ListDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
) error {
	return p.RunForProject(ctx, options.ProjectInfo, append([]string{"list"}, options.ExtraArgs...))
}

func (p *Pip) AuditDependencies(
	_ context.Context,
	_ internal.ProjectDependencyManagerOptions,
) error {
	return errors.New("not supported")
}
