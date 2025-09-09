package internal

// A Context that Tools execute on
type Context struct {
	// Config is an initial project configuration.
	Config Config

	// Project describes a project
	Project Project

	// Workflow is a workflow to use
	Workflow Workflow

	// Tools store all supported tools
	Tools Tools

	// EnvironmentProviders stores all supported environment providers
	EnvironmentProviders EnvironmentProviders

	// Environment is where the tool exists and runs
	Environment Environment
}
