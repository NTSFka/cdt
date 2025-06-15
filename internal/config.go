package internal

type Config struct {
	// RootDirectory is a root directory of the project
	RootDirectory string

	// BuildDirectory is the project's build directory
	BuildDirectory *string
}
