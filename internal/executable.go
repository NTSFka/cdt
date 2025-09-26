package internal

import (
	"context"
	"io"
	"strings"
)

// RunOptions are options for executing tool.
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

// ExecutableRuntime represents the runtime in which the executable runs.
type ExecutableRuntime interface {
	// The Id returns the runtime unique identifier
	Id() string

	// RunExecutable runs the executable
	RunExecutable(ctx context.Context, options RunOptions, path string, args []string) error
}

// Executable represent an executable.
type Executable struct {
	// Path to executable
	Path string

	// Args stores additional arguments that will be used always
	Args []string

	// Runtime in which the executable will be run
	Runtime ExecutableRuntime
}

// Run starts the executable with the given arguments.
func (e *Executable) Run(ctx context.Context, options RunOptions, args []string) error {
	runArgs := e.buildArgs(args)

	if !options.Silent {
		Info("%v: %v %v", e.Runtime.Id(), e.Path, strings.Join(runArgs, " "))
	}

	return Trace(
		ctx,
		"executable.run",
		func() error {
			return e.Runtime.RunExecutable(ctx, options, e.Path, runArgs)
		},
		"path", e.Path,
		"args", args,
		"directory", options.Directory,
	)
}

func (e *Executable) buildArgs(args []string) []string {
	eArgs := e.Args

	if eArgs == nil {
		eArgs = []string{}
	}

	return append(eArgs, args...)
}
