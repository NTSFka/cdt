package tool

import "cdt/internal"

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

func (c *Composer) AddDependencies(project internal.Project, dependencies []string, dev bool) error {
	args := []string{"require"}

	if dev {
		args = append(args, "--dev")
	}

	return c.RunForProject(project, append(args, dependencies...))
}

func (c *Composer) RemoveDependencies(project internal.Project, dependencies []string, dev bool) error {
	args := []string{"remove"}

	if dev {
		args = append(args, "--dev")
	}

	return c.RunForProject(project, append(args, dependencies...))
}

func (c *Composer) UpdateDependencies(project internal.Project, dependencies []string) error {
	return c.RunForProject(project, append([]string{"update"}, dependencies...))
}

func (c *Composer) FetchDependencies(project internal.Project, noDev bool) error {
	args := []string{"install"}

	if noDev {
		args = append(args, "--no-dev")
	}

	return c.RunForProject(project, args)
}

func (c *Composer) ListDependencies(project internal.Project) error {
	return c.RunForProject(project, []string{"show"})
}

func (c *Composer) AuditDependencies(project internal.Project) error {
	return c.RunForProject(project, []string{"audit"})
}
