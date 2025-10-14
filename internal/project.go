package internal

import (
	"context"
	"errors"
	"slices"
)

// A ProjectStructureProvider provides a detailed structure of the project.
type ProjectStructureProvider interface {
	// Structure returns project structure
	Structure(ctx context.Context, info ProjectInfo) (*ProjectStructure, error)
}

// ProjectStructure describes a project structure.
type ProjectStructure struct {
	Targets map[string]ProjectTarget
}

// ProjectTarget describes a project target.
type ProjectTarget struct {
	Files      []string
	Dependency bool
}

// GetFiles returns all files in the project.
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

var ErrNoIntermediateDirectory = errors.New("no intermediate directory set")

// A ProjectInfo provides information about the project.
type ProjectInfo struct {
	// Directory is the directory of the project
	Directory string

	// IntermediateDirectory is the directory of the project's intermediate files
	IntermediateDirectory *string

	// StructureProvider provides structure of the project primary in cases when the structure is dynamic from
	// some configuration
	StructureProvider ProjectStructureProvider
}

// Structure returns project structure.
func (p *ProjectInfo) Structure(ctx context.Context) (*ProjectStructure, error) {
	return p.StructureProvider.Structure(ctx, *p)
}

// A EmptyProjectStructureProvider provides detailed empty project structure.
type EmptyProjectStructureProvider struct {
}

// Structure returns project structure.
func (p *EmptyProjectStructureProvider) Structure(
	_ context.Context,
	_ ProjectInfo,
) (*ProjectStructure, error) {
	return nil, errors.New("no structure provided")
}

// A FixedProjectStructureProvider provides a predefined project structure.
type FixedProjectStructureProvider struct {
	ProjectStructure ProjectStructure
}

// Structure returns project structure.
func (p *FixedProjectStructureProvider) Structure(
	_ context.Context,
	_ ProjectInfo,
) (*ProjectStructure, error) {
	return &p.ProjectStructure, nil
}

// A Project describes a project and its workflow.
type Project struct {
	// Info provides information about the project
	Info ProjectInfo

	// Workflow describes a project workflow
	Workflow Workflow
}
