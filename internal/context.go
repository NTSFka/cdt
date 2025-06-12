package internal

// A Context that Tools execute on
type Context struct {
	// Project is a project to work with
	Project Project

	// A Workflow that will be used
	Workflow Workflow

	// Tools store all supported tools
	Tools Tools
}
