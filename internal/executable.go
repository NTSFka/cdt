package internal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// RunContext is a context in which the executable will run.
type RunContext struct {
	Directory string
	Output    io.Writer
	Error     io.Writer
}

// NewRunContext create context for directory and standard output
func NewRunContext(directory string) RunContext {
	return RunContext{Directory: directory, Output: os.Stdout, Error: os.Stderr}
}

// Executable is a structure that stores information about executable
type Executable struct {
	// Path to executable
	Path string

	// RunFunc is a function that will run the executable
	RunFunc func(ctx RunContext, path string, args []string) error
}

// Run starts the executable with the given arguments
func (t *Executable) Run(ctx RunContext, args []string) error {
	fmt.Printf("RUN: %s %v\n", t.Path, strings.Join(args, " "))

	return t.RunFunc(ctx, t.Path, args)
}
