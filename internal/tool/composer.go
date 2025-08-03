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
			detect,
		),
	}
}

func (c *Composer) AddDependencies(project internal.Project, dependencies []string) error {
	return c.RunForProject(project, append([]string{"require"}, dependencies...))
}

func (c *Composer) RemoveDependencies(project internal.Project, dependencies []string) error {
	return c.RunForProject(project, append([]string{"remove"}, dependencies...))
}

func (c *Composer) UpdateDependencies(project internal.Project, dependencies []string) error {
	return c.RunForProject(project, append([]string{"update"}, dependencies...))
}

func (c *Composer) FetchDependencies(project internal.Project) error {
	return c.RunForProject(project, []string{"install"})
}

func (c *Composer) ListDependencies(project internal.Project) error {
	return c.RunForProject(project, []string{"show"})
}

func (c *Composer) AuditDependencies(project internal.Project) error {
	return c.RunForProject(project, []string{"audit"})
}
