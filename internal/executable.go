package internal

import (
	"fmt"
)

// Executable is a structure that stores information about executable
type Executable struct {
	// Path to executable
	Path string

	// Environment in which the executable exists
	Environment Environment
}

// Run starts the executable with the given arguments
func (t *Executable) Run(args []string) error {
	fmt.Printf("RUN: %s\n", args)

	return t.Environment.RunExecutable(t.Path, args)
}
