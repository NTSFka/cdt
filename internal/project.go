package internal

import "slices"

// A ProjectStructureProvider provides a detailed structure of the project
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
	Directory             string
	IntermediateDirectory *string
	StructureProvider     ProjectStructureProvider
}

// RootDirectory returns the project root directory
func (p *Project) RootDirectory() string {
	return p.Directory
}

// BuildDirectory returns the project build directory
func (p *Project) BuildDirectory() string {
	if p.IntermediateDirectory == nil {
		panic("intermediate directory is not set")
	}

	return *p.IntermediateDirectory
}

// Structure returns project structure
func (p *Project) Structure() (*ProjectStructure, error) {
	return p.StructureProvider.Structure(*p)
}

// A EmptyProjectStructureProvider provides detailed empty project structure
type EmptyProjectStructureProvider struct {
}

// Structure returns project structure
func (p *EmptyProjectStructureProvider) Structure(_ Project) (*ProjectStructure, error) {
	return &ProjectStructure{}, nil
}

// A FixedProjectStructureProvider provides a predefined project structure
type FixedProjectStructureProvider struct {
	ProjectStructure ProjectStructure
}

// Structure returns project structure
func (p *FixedProjectStructureProvider) Structure(_ Project) (*ProjectStructure, error) {
	return &p.ProjectStructure, nil
}
