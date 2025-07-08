package internal

// A Context that Tools execute on
type Context struct {
	// Config is an initial project configuration.
	Config Config

	// Project is a project to work with
	Project Project

	// Tools store all supported tools
	Tools Tools

	// EnvironmentProviders stores all supported environment providers
	EnvironmentProviders EnvironmentProviders

	// Environment is where the tool exists and runs
	Environment Environment
}
