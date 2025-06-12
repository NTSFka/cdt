package internal

import "slices"

// ProjectStructure describes a project structure
type ProjectStructure struct {
	Targets map[string]ProjectTarget
}

// ProjectTarget describes a project target
type ProjectTarget struct {
	Files      []string
	Dependency bool
}

// GetFiles returns all files in the project
func (p *ProjectStructure) GetFiles() []string {
	var files []string

	for _, target := range p.Targets {
		if !target.Dependency {
			files = append(files, target.Files...)
		}
	}

	slices.Sort(files)
	return slices.Compact(files)
}

// A Project describes a project in a specific directory
type Project interface {
	RootDirectory() string
	BuildDirectory() string

	// Structure returns project structure
	Structure() (*ProjectStructure, error)
}
