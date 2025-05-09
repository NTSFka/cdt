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

// A ProjectConfigurator allow configuring a project
type ProjectConfigurator interface {
	// Configure the project
	Configure() error
}

// A ProjectBuilder allow building a project
type ProjectBuilder interface {
	// BuildAll builds all targets in the project
	BuildAll() error

	// BuildTarget builds a specific target in the project
	BuildTarget(target string) error

	// BuildTargets builds specific targets in the project
	BuildTargets(targets []string) error
}

// A ProjectTester allow testing a project
type ProjectTester interface {
	// TestAll runs all tests in the project
	TestAll() error

	// Test runs tests that match the pattern
	Test(pattern string) error
}

// A ProjectFormatter allow formatting files of a project
type ProjectFormatter interface {
	FormatAll() error
	FormatFiles(filenames []string) error
	FormatCheckAll() error
	FormatCheckFiles(filenames []string) error
}

// A ProjectLinter allow linting files of a project
type ProjectLinter interface {
	LintAll() error
	LintFiles(filenames []string) error
}

// A ProjectRunner allow running executables of a project
type ProjectRunner interface {
	// Run a target
	Run(target string, args []string) error
}

// A Project describes a project in a specific directory
type Project interface {
	RootDirectory() string
	BuildDirectory() string

	// Structure returns project structure
	Structure() (*ProjectStructure, error)

	// Configurator returns a configurator for the project
	Configurator() ProjectConfigurator

	// Builder returns a builder for the project
	Builder() ProjectBuilder

	// Tester returns a tester for the project
	Tester() ProjectTester

	// Formatter returns a formatter for the project
	Formatter() ProjectFormatter

	// Linter returns a linter for the project
	Linter() ProjectLinter

	// Runner returns a runner of the project
	Runner() ProjectRunner
}
