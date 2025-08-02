package internal

// A ProjectConfigurator allow configuring a project
type ProjectConfigurator interface {
	// Configure the project
	Configure(project Project, args []string) error
}

// A ProjectBuilder allow building a project
type ProjectBuilder interface {
	// BuildAll builds all targets in the project
	BuildAll(project Project, args []string) error

	// BuildTargets builds specific targets in the project
	BuildTargets(project Project, targets []string, args []string) error
}

// A ProjectTester allow testing a project
type ProjectTester interface {
	// TestAll runs all tests in the project
	TestAll(project Project, args []string) error

	// Test runs tests that match the pattern
	Test(project Project, pattern string, args []string) error
}

// A ProjectFormatter allow formatting files of a project
type ProjectFormatter interface {
	// FormatAll formates all files in the project
	FormatAll(project Project, args []string) error

	// FormatFiles formates specified files in the project
	FormatFiles(project Project, filenames []string, args []string) error

	// FormatCheckAll check if all files in the project are formatted
	FormatCheckAll(project Project, args []string) error

	// FormatCheckFiles check if all specified files in the project are formatted
	FormatCheckFiles(project Project, filenames []string, args []string) error
}

// A ProjectLinter allow linting files of a project
type ProjectLinter interface {
	// LintAll lints all project files
	LintAll(project Project, args []string) error

	// LintFiles perform linting on specified files
	LintFiles(project Project, filenames []string, args []string) error
}

// A ProjectRunner allow running executables of a project
type ProjectRunner interface {
	// RunTarget run a target
	RunTarget(project Project, target string, args []string) error
}

// ProjectDependencyManager manages project dependencies (libraries, packaged, etc.)
type ProjectDependencyManager interface {
	// AddDependencies adds new dependencies to the project
	AddDependencies(project Project, dependency ...string) error

	// RemoveDependencies removes the dependencies from the project
	RemoveDependencies(project Project, dependencies ...string) error

	// UpdateDependencies updates specified dependency in the project (empty dependency list means update all)
	UpdateDependencies(project Project, dependency ...string) error

	// FetchDependencies fetches all specified dependencies to the project
	FetchDependencies(project Project) error

	// ListDependencies lists all specified dependencies in the project
	ListDependencies(project Project) error

	// AuditDependencies audits all specified dependencies in the project for security issues
	AuditDependencies(project Project) error
}

// A Workflow describes how to work on a project
type Workflow struct {
	// Configurator stores a configurator for the project
	Configurator ProjectConfigurator

	// Builder stores a builder for the project
	Builder ProjectBuilder

	// Tester stores a tester for the project
	Tester ProjectTester

	// Formatter stores a formatter for the project
	Formatter ProjectFormatter

	// Linter stores a linter for the project
	Linter ProjectLinter

	// Runner stores a runner of the project
	Runner ProjectRunner

	// DependencyManager stores a dependency manager of the project
	DependencyManager ProjectDependencyManager
}
