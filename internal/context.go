package internal

// A Context that Tools execute on
type Context struct {
	// Project is a project to work with
	Project Project

	// Tools store all supported tools
	Tools Tools
}
