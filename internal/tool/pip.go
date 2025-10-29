package tool

import (
	"context"
	"errors"

	"cdt/internal"
)

const IdPip = "pip"

type Pip struct {
	internal.ExecutableTool
}

// DetectPip create a tool for pip.
func DetectPip(
	ctx context.Context,
	options DetectOptions,
) *Pip {
	if path, ok := options.ToolsPaths[IdPip]; ok {
		return NewPip(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewPip(internal.DetectExecutableChain(
		[]string{"pip", "pip3"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewPip creates a pip tool from a custom executable.
func NewPip(detect internal.ExecutableToolDetectFunc) *Pip {
	return &Pip{
		ExecutableTool: internal.MakeExecutableTool(
			IdPip,
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
