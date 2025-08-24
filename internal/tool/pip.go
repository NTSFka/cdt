package tool

import (
	"cdt/internal"
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
			detect,
		),
	}
}

func (p *Pip) AddDependencies(project internal.Project, dependencies []string, _ bool) error {
	return p.RunForProject(project, append([]string{"install"}, dependencies...))
}

func (p *Pip) RemoveDependencies(project internal.Project, dependencies []string, _ bool) error {
	return p.RunForProject(project, append([]string{"uninstall"}, dependencies...))
}

func (p *Pip) UpdateDependencies(project internal.Project, dependencies []string) error {
	return p.RunForProject(project, append([]string{"install", "--upgrade"}, dependencies...))
}

func (p *Pip) FetchDependencies(project internal.Project, _ bool) error {
	return p.RunForProject(project, []string{"install", "-r", "requirements.txt"})
}

func (p *Pip) ListDependencies(project internal.Project) error {
	return p.RunForProject(project, []string{"list"})
}

func (p *Pip) AuditDependencies(_ internal.Project) error {
	return errors.New("not supported")
}
