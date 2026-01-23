package internal

import "context"

// OutputOptions specifies output options for various operations.
type OutputOptions[T any] struct {
	Format   T
	Filename *string
}

// ProjectConfiguratorOptions are options for configuring a project.
type ProjectConfiguratorOptions struct {
	ProjectInfo

	// ExtraArgs are extra arguments for the specific configurator implementation
	ExtraArgs []string
}

// A ProjectConfigurator allow configuring a project.
type ProjectConfigurator interface {
	// Configure the project
	Configure(ctx context.Context, options ProjectConfiguratorOptions) error
}

// ProjectBuilderOptions are options for building a project.
type ProjectBuilderOptions struct {
	ProjectInfo

	// ExtraArgs are extra arguments for the specific builder implementation
	ExtraArgs []string
}

// A ProjectBuilder allow building a project.
type ProjectBuilder interface {
	// BuildAll builds all targets in the project
	BuildAll(ctx context.Context, options ProjectBuilderOptions) error

	// BuildTargets builds specific targets in the project
	BuildTargets(ctx context.Context, options ProjectBuilderOptions, targets []string) error
}

// TestsReportFormat specifies the format of the tests report.
type TestsReportFormat string

const (
	TestsReportFormatRaw  TestsReportFormat = "raw"
	TestsReportFormatJson TestsReportFormat = "json"
	TestsReportFormatCtrf TestsReportFormat = "ctrf"
)

// ProjectTesterOptions are options for testing a project.
type ProjectTesterOptions struct {
	ProjectInfo

	// ExtraArgs are extra arguments for the specific tester implementation
	ExtraArgs []string

	// Output specifies how the operation should produce output
	Output OutputOptions[TestsReportFormat]
}

// A ProjectTester allow testing a project.
type ProjectTester interface {
	// TestAll runs all tests in the project
	TestAll(ctx context.Context, options ProjectTesterOptions) error

	// TestPattern runs tests that match the pattern
	TestPattern(ctx context.Context, options ProjectTesterOptions, pattern string) error
}

// ProjectCoverageCollectorOptions are options for collecting coverage of a project.
type ProjectCoverageCollectorOptions struct {
	ProjectInfo

	// ExtraArgs are extra arguments for the specific collector implementation
	ExtraArgs []string
}

// A ProjectCoverageCollector allow coverage collecting for a project.
type ProjectCoverageCollector interface {
	// CollectCoverageAll runs all tests in the project and collects coverage
	CollectCoverageAll(ctx context.Context, options ProjectCoverageCollectorOptions) error

	// CollectCoveragePattern runs tests that match the pattern and collects coverage
	CollectCoveragePattern(
		ctx context.Context,
		options ProjectCoverageCollectorOptions,
		pattern string,
	) error
}

// ProjectFormatterOptions are options for formatting a project.
type ProjectFormatterOptions struct {
	ProjectInfo

	// ExtraArgs are extra arguments for the specific formatter implementation
	ExtraArgs []string
}

// A ProjectFormatter allow formatting files of a project.
type ProjectFormatter interface {
	// FormatAll formates all files in the project
	FormatAll(ctx context.Context, options ProjectFormatterOptions) error

	// FormatFiles formates specified files in the project
	FormatFiles(ctx context.Context, options ProjectFormatterOptions, filenames []string) error

	// FormatCheckAll check if all files in the project are formatted
	FormatCheckAll(ctx context.Context, options ProjectFormatterOptions) error

	// FormatCheckFiles check if all specified files in the project are formatted
	FormatCheckFiles(ctx context.Context, options ProjectFormatterOptions, filenames []string) error
}

// ProjectLinterOptions are options for linting a project.
type ProjectLinterOptions struct {
	ProjectInfo

	// ExtraArgs are extra arguments for the specific linter implementation
	ExtraArgs []string
}

// A ProjectLinter allow linting files of a project.
type ProjectLinter interface {
	// LintAll lints all project files
	LintAll(ctx context.Context, options ProjectLinterOptions) error

	// LintFiles perform linting on specified files
	LintFiles(ctx context.Context, options ProjectLinterOptions, filenames []string) error
}

// ProjectRunnerOptions are options for running a target in the project.
type ProjectRunnerOptions struct {
	ProjectInfo

	// Runtime is a runtime of the target binary
	Runtime ExecutableRuntime

	// ExtraArgs are extra arguments for the specific runner implementation
	ExtraArgs []string
}

// A ProjectRunner allow running executables of a project.
type ProjectRunner interface {
	// RunTarget run a target
	RunTarget(ctx context.Context, options ProjectRunnerOptions, target string) error
}

// ProjectDependencyManagerOptions are options for managing of project dependencies.
type ProjectDependencyManagerOptions struct {
	ProjectInfo

	// ExtraArgs are extra arguments for the specific dependency manager implementation
	ExtraArgs []string
}

// ProjectDependencyManager manages project dependencies (libraries, packaged, etc.)
type ProjectDependencyManager interface {
	// AddDependencies adds new dependencies to the project
	AddDependencies(
		ctx context.Context,
		options ProjectDependencyManagerOptions,
		dependencies []string,
		dev bool,
	) error

	// RemoveDependencies removes the dependencies from the project
	RemoveDependencies(
		ctx context.Context,
		options ProjectDependencyManagerOptions,
		dependencies []string,
		dev bool,
	) error

	// UpdateDependencies updates specified dependencies in the project (empty dependencies mean update all)
	UpdateDependencies(
		ctx context.Context,
		options ProjectDependencyManagerOptions,
		dependencies []string,
	) error

	// FetchDependencies fetches all specified dependencies to the project
	FetchDependencies(
		ctx context.Context,
		options ProjectDependencyManagerOptions,
		noDev bool,
	) error

	// ListDependencies lists all specified dependencies in the project
	ListDependencies(ctx context.Context, options ProjectDependencyManagerOptions) error

	// AuditDependencies audits all specified dependencies in the project for security issues
	AuditDependencies(ctx context.Context, options ProjectDependencyManagerOptions) error
}

// A Workflow describes how to work on a project.
type Workflow struct {
	// Name is the name of the workflow
	Name string

	// Configurator stores a configurator for the project
	Configurator ProjectConfigurator

	// Builder stores a builder for the project
	Builder ProjectBuilder

	// Tester stores a tester for the project
	Tester ProjectTester

	// CoverageCollector stores a coverage collector for the project
	CoverageCollector ProjectCoverageCollector

	// Formatter stores a formatter for the project
	Formatter ProjectFormatter

	// Linter stores a linter for the project
	Linter ProjectLinter

	// Runner stores a runner of the project
	Runner ProjectRunner

	// DependencyManager stores a dependency manager of the project
	DependencyManager ProjectDependencyManager
}
