package internal

import (
	"context"
	"io"
)

// RunOptions are options for executing tool
type RunOptions struct {
	Directory string
	Input     io.Reader
	Output    io.Writer
	Error     io.Writer
}

// Executable is a structure that stores information about executable
type Executable struct {
	// Path to executable
	Path string

	// RunFunc is a function that will run the executable
	RunFunc func(ctx context.Context, options RunOptions, path string, args []string) error
}

// Run starts the executable with the given arguments
func (t *Executable) Run(ctx context.Context, options RunOptions, args []string) error {
	return Trace(
		ctx,
		"executable.run",
		func() error {
			return t.RunFunc(ctx, options, t.Path, args)
		},
		"path", t.Path,
		"args", args,
		"directory", options.Directory,
	)
}
