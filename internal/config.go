package internal

type ConfigEnvironment struct {
	ToolName string
	Argument string
}

type Config struct {
	// RootDirectory is a root directory of the project
	RootDirectory string

	// BuildDirectory is the project's build directory
	BuildDirectory *string

	// Environment defines an environment to use
	Environment *ConfigEnvironment
}
