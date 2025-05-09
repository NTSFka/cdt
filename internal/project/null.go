package project

import (
	. "cdt/internal"
)

// A NullProject describes a project that cannot be used
type NullProject struct {
	Project

	rootDirectory  string
	buildDirectory string
}

func (p *NullProject) Configurator() ProjectConfigurator {
	return nil
}

func (p *NullProject) Builder() ProjectBuilder {
	return nil
}

func (p *NullProject) Tester() ProjectTester {
	return nil
}

func (p *NullProject) Formatter() ProjectFormatter {
	return nil
}

func (p *NullProject) Linter() ProjectLinter {
	return nil
}

func (p *NullProject) Runner() ProjectRunner { return nil }

// NewNullProject creates a new null project
func NewNullProject(directory string, buildDirectory string) NullProject {
	return NullProject{
		rootDirectory:  directory,
		buildDirectory: buildDirectory,
	}
}

func (p *NullProject) RootDirectory() string {
	return p.rootDirectory
}

func (p *NullProject) BuildDirectory() string {
	return p.buildDirectory
}

func (p *NullProject) Structure() (*ProjectStructure, error) {
	info := ProjectStructure{
		Targets: make(map[string]ProjectTarget),
	}

	return &info, nil
}
