package internal

import "slices"

// A ProjectStructureProvider provides detailed structure of the project
type ProjectStructureProvider interface {
	// Structure returns project structure
	Structure(project Project) (*ProjectStructure, error)
}

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
type Project struct {
	rootDirectory     string
	buildDirectory    string
	structureProvider ProjectStructureProvider
}

// MakeProject creates a new project
func MakeProject(rootDirectory string, buildDirectory string, structureProvider ProjectStructureProvider) Project {
	return Project{
		rootDirectory:     rootDirectory,
		buildDirectory:    buildDirectory,
		structureProvider: structureProvider,
	}
}

// RootDirectory returns project root directory
func (p *Project) RootDirectory() string {
	return p.rootDirectory
}

// BuildDirectory returns project build directory
func (p *Project) BuildDirectory() string {
	return p.buildDirectory
}

// Structure returns project structure
func (p *Project) Structure() (*ProjectStructure, error) {
	return p.structureProvider.Structure(*p)
}

// A EmptyProjectStructureProvider provides detailed empty project structure
type EmptyProjectStructureProvider struct {
}

// Structure returns project structure
func (p *EmptyProjectStructureProvider) Structure(project Project) (*ProjectStructure, error) {
	return &ProjectStructure{}, nil
}
