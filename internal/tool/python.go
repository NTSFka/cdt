package tool

import (
	"cdt/internal"
	"context"
	"errors"
	"os/exec"
	"path/filepath"
)

type Python struct {
	internal.ExecutableTool
}

// DetectPython create a tool for python
func DetectPython(ctx context.Context, environment internal.Environment) *Python {
	return NewPython(func() *internal.Executable {
		return environment.FindExecutable(ctx, "python3")
	})
}

// NewPython creates a python tool from a custom executable
func NewPython(detect func() *internal.Executable) *Python {
	return &Python{
		ExecutableTool: internal.MakeExecutableTool(
			"python",
			"Python",
			"Python is a programming language that lets you work quickly and integrate systems more effectively.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagRun, internal.ToolTagEnvironment},
			detect,
		),
	}
}

func (p *Python) RunTarget(ctx context.Context, options internal.ProjectRunnerOptions, target string) error {
	return p.RunForProject(ctx, options.ProjectInfo, append([]string{target}, options.ExtraArgs...))
}

func (p *Python) Aliases() []string {
	return []string{"pyenv"}
}

func (p *Python) ParameterInfo() string {
	return "path to virtual environment directory"
}

func (p *Python) Detect(directory string) *internal.Environment {
	// Check if the directory contains a venv
	for _, dir := range []string{"venv", ".venv"} {
		if internal.PathExists(filepath.Join(directory, dir)) {
			env, _ := p.CreateEnvironment(directory, dir)
			return &env
		}
	}

	return nil
}

// CreateEnvironment create python virtual environment where the service is used for running tools
func (p *Python) CreateEnvironment(directory, path string) (internal.Environment, error) {
	internal.Debug("pyenv.create_environment", "directory", directory, "path", path)

	if len(path) == 0 {
		return nil, errors.New("python virtual environment path is required")
	}

	env := pythonVirtualEnvironment{
		directory:     directory,
		venvDirectory: path,
		python:        p,
	}

	return &env, nil
}

type pythonVirtualEnvironment struct {
	directory     string
	venvDirectory string
	python        *Python
}

func (e *pythonVirtualEnvironment) Id() string {
	return "pyenv"
}

func (e *pythonVirtualEnvironment) Start(ctx context.Context) error {
	// Check if the environment already exists
	if !internal.PathExists(filepath.Join(e.venvDirectory, "pyvenv.cfg")) {
		options := internal.RunOptions{
			Directory: e.directory,
			Input:     nil,
			Output:    nil,
			Error:     nil,
			Silent:    true,
		}

		if err := e.python.Run(ctx, options, []string{"-m", "venv", e.venvDirectory}); err != nil {
			return err
		}
	}

	return nil
}

func (e *pythonVirtualEnvironment) IsRunning(_ context.Context) bool {
	// Environment is always running
	return true
}

func (e *pythonVirtualEnvironment) Stop(_ context.Context) error {
	// Environment is always running
	return nil
}

func (e *pythonVirtualEnvironment) Cleanup(_ context.Context) error {
	// Environment is always running
	return nil
}

func (e *pythonVirtualEnvironment) findPath(name string) *string {
	if path := filepath.Join(e.venvDirectory, "bin", name); internal.PathExists(path) {
		return &path
	}

	// Windows
	if path := filepath.Join(e.venvDirectory, "Scripts", name); internal.PathExists(path) {
		return &path
	}

	return nil
}

func (e *pythonVirtualEnvironment) FindExecutable(ctx context.Context, name string) *internal.Executable {
	return internal.Trace(ctx, "pyenv.find_executable", func() *internal.Executable {
		if path := e.findPath(name); path != nil {
			return &internal.Executable{
				Path:    name,
				Runtime: e,
			}
		}

		return nil
	}, "venv", e.venvDirectory, "name", name)
}

func (e *pythonVirtualEnvironment) RunExecutable(ctx context.Context, options internal.RunOptions, name string, args []string) error {
	path := e.findPath(name)

	if path == nil {
		return errors.New("executable not found")
	}

	command := exec.CommandContext(ctx, *path, args...)
	command.Dir = options.Directory
	command.Stdin = options.Input
	command.Stdout = options.Output
	command.Stderr = options.Error

	return internal.Trace(ctx, "pyenv.run", func() error {
		return command.Run()
	}, "venv", e.venvDirectory, "path", *path, "args", args)
}
