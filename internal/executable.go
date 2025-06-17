package internal

import (
	"fmt"
)

// Executable is a structure that stores information about executable
type Executable struct {
	// Path to executable
	Path string

	// RunFunc is a function that will run the executable
	RunFunc func(path string, args []string) error
}

// Run starts the executable with the given arguments
func (t *Executable) Run(args []string) error {
	fmt.Printf("RUN: %s\n", args)

	return t.RunFunc(t.Path, args)
}
