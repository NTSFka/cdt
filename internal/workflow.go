package internal

// A ProjectConfigurator allow configuring a project
type ProjectConfigurator interface {
	// Configure the project
	Configure(project Project) error
}

// A ProjectBuilder allow building a project
type ProjectBuilder interface {
	// BuildAll builds all targets in the project
	BuildAll(project Project) error

	// BuildTargets builds specific targets in the project
	BuildTargets(project Project, targets []string) error
}

// A ProjectTester allow testing a project
type ProjectTester interface {
	// TestAll runs all tests in the project
	TestAll(project Project) error

	// Test runs tests that match the pattern
	Test(project Project, pattern string) error
}

// A ProjectFormatter allow formatting files of a project
type ProjectFormatter interface {
	FormatAll(project Project) error
	FormatFiles(project Project, filenames []string) error
	FormatCheckAll(project Project) error
	FormatCheckFiles(project Project, filenames []string) error
}

// A ProjectLinter allow linting files of a project
type ProjectLinter interface {
	LintAll(project Project) error
	LintFiles(project Project, filenames []string) error
}

// A ProjectRunner allow running executables of a project
type ProjectRunner interface {
	// Run a target
	Run(project Project, target string, args []string) error
}

// A Workflow describes how to work on a project
type Workflow struct {
	// Configurator returns a configurator for the project
	Configurator ProjectConfigurator

	// Builder returns a builder for the project
	Builder ProjectBuilder

	// Tester returns a tester for the project
	Tester ProjectTester

	// Formatter returns a formatter for the project
	Formatter ProjectFormatter

	// Linter returns a linter for the project
	Linter ProjectLinter

	// Runner returns a runner of the project
	Runner ProjectRunner
}
