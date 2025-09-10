package tool

import (
	"cdt/internal"
	"context"
)

type Composer struct {
	internal.ExecutableTool
}

// DetectComposer create a tool for composer
func DetectComposer(environment internal.Environment) *Composer {
	return NewComposer(func() *internal.Executable {
		// PHAR
		if executable := environment.FindExecutable("composer.phar"); executable != nil {
			return executable
		}

		// System version
		if executable := environment.FindExecutable("composer"); executable != nil {
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

func (c *Composer) AddDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string, dev bool) error {
	args := []string{"require"}

	if dev {
		args = append(args, "--dev")
	}

	return c.RunForProject(ctx, info, append(args, dependencies...))
}

func (c *Composer) RemoveDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string, dev bool) error {
	args := []string{"remove"}

	if dev {
		args = append(args, "--dev")
	}

	return c.RunForProject(ctx, info, append(args, dependencies...))
}

func (c *Composer) UpdateDependencies(ctx context.Context, info internal.ProjectInfo, dependencies []string) error {
	return c.RunForProject(ctx, info, append([]string{"update"}, dependencies...))
}

func (c *Composer) FetchDependencies(ctx context.Context, info internal.ProjectInfo, noDev bool) error {
	args := []string{"install"}

	if noDev {
		args = append(args, "--no-dev")
	}

	return c.RunForProject(ctx, info, args)
}

func (c *Composer) ListDependencies(ctx context.Context, info internal.ProjectInfo) error {
	return c.RunForProject(ctx, info, []string{"show"})
}

func (c *Composer) AuditDependencies(ctx context.Context, info internal.ProjectInfo) error {
	return c.RunForProject(ctx, info, []string{"audit"})
}
