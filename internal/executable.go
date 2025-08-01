package internal

import (
	"context"
	"io"
	"strings"
)

// RunOptions are options for executing tool
type RunOptions struct {
	// Directory in which executable should be run
	Directory string
	// Input provides executable input
	Input io.Reader
	// Output captures executable output
	Output io.Writer
	// Error captures executable error output
	Error io.Writer
	// Silent is a flag that disables cdt run info
	Silent bool
}

// Executable represent an executable
type Executable struct {
	// Path to executable
	Path string

	// Args stores additional arguments that will be used always
	Args []string

	// RunFunc is a function that will run the executable
	RunFunc func(ctx context.Context, options RunOptions, path string, args []string) error
}

func (e *Executable) buildArgs(args []string) []string {
	eArgs := e.Args

	if eArgs == nil {
		eArgs = []string{}
	}

	return append(eArgs, args...)
}

// Run starts the executable with the given arguments
func (e *Executable) Run(ctx context.Context, options RunOptions, args []string) error {
	runArgs := e.buildArgs(args)

	if !options.Silent {
		Info("%v %v", e.Path, strings.Join(runArgs, " "))
	}

	return Trace(
		ctx,
		"executable.run",
		func() error {
			return e.RunFunc(ctx, options, e.Path, runArgs)
		},
		"path", e.Path,
		"args", args,
		"directory", options.Directory,
	)
}
