package tool

import (
	"cdt/internal"
	"context"
)

type Composer struct {
	internal.ExecutableTool
}

// DetectComposer create a tool for composer
func DetectComposer(ctx context.Context, environment internal.Environment) *Composer {
	return NewComposer(func() *internal.Executable {
		// PHAR
		if executable := environment.FindExecutable(ctx, "composer.phar"); executable != nil {
			return executable
		}

		// System version
		if executable := environment.FindExecutable(ctx, "composer"); executable != nil {
			return executable
		}

		return nil
	})
}

// NewComposer creates a composer tool from a custom executable
func NewComposer(detect func() *internal.Executable) *Composer {
	return &Composer{
		ExecutableTool: internal.MakeExecutableTool(
			"composer",
			"Composer",
			"A Dependency Manager for PHP",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagDependency},
			detect,
		),
	}
}

func (c *Composer) AddDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string, dev bool) error {
	args := append([]string{"require"}, options.ExtraArgs...)

	if dev {
		args = append(args, "--dev")
	}

	return c.RunForProject(ctx, options.ProjectInfo, append(args, dependencies...))
}

func (c *Composer) RemoveDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string, dev bool) error {
	args := append([]string{"remove"}, options.ExtraArgs...)

	if dev {
		args = append(args, "--dev")
	}

	return c.RunForProject(ctx, options.ProjectInfo, append(args, dependencies...))
}

func (c *Composer) UpdateDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, dependencies []string) error {
	return c.RunForProject(ctx, options.ProjectInfo, append(append([]string{"update"}, options.ExtraArgs...), dependencies...))
}

func (c *Composer) FetchDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions, noDev bool) error {
	args := append([]string{"install"}, options.ExtraArgs...)

	if noDev {
		args = append(args, "--no-dev")
	}

	return c.RunForProject(ctx, options.ProjectInfo, args)
}

func (c *Composer) ListDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions) error {
	return c.RunForProject(ctx, options.ProjectInfo, append([]string{"show"}, options.ExtraArgs...))
}

func (c *Composer) AuditDependencies(ctx context.Context, options internal.ProjectDependencyManagerOptions) error {
	return c.RunForProject(ctx, options.ProjectInfo, append([]string{"audit"}, options.ExtraArgs...))
}
