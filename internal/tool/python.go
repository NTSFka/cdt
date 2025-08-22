package tool

import (
	"cdt/internal"
	"context"
	"os/exec"
	"path/filepath"
)

type Python struct {
	internal.ExecutableTool
}

// DetectPython create a tool for python
func DetectPython(environment internal.Environment) *Python {
	return NewPython(func() *internal.Executable {
		return environment.FindExecutable("python3")
	})
}

// NewPython creates a python tool from a custom executable
func NewPython(detect func() *internal.Executable) *Python {
	return &Python{
		ExecutableTool: internal.MakeExecutableTool(
			"python",
			"Python",
			"Python is a programming language that lets you work quickly and integrate systems more effectively.",
			detect,
		),
	}
}

func (p *Python) RunTarget(project internal.Project, target string, args []string) error {
	return p.RunForProject(project, append([]string{target}, args...))
}

func (p *Python) IdShort() string {
	return "pyenv"
}

// CreateEnvironment create python virtual environment where the service is used for running tools
func (p *Python) CreateEnvironment(directory, path string) (internal.Environment, error) {
	internal.Debug("pyenv.create_environment", "directory", directory, "path", path)

	env := pythonVirtualEnvironment{
		directory:     directory,
		venvDirectory: path,
	}

	return &env, nil
}

type pythonVirtualEnvironment struct {
	directory     string
	venvDirectory string
}

func (e *pythonVirtualEnvironment) Id() string {
	return "pyenv"
}

func (e *pythonVirtualEnvironment) Start(_ context.Context) error {
	// Environment is always running
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

func (e *pythonVirtualEnvironment) FindExecutable(name string) *internal.Executable {
	ctx := context.Background()

	return internal.Trace(ctx, "pyenv.find_executable", func() *internal.Executable {
		if path := e.findPath(name); path != nil {
			return &internal.Executable{
				Path:    *path,
				Runtime: e,
			}
		}

		return nil
	}, "venv", e.venvDirectory, "name", name)
}

func (e *pythonVirtualEnvironment) RunExecutable(ctx context.Context, options internal.RunOptions, path string, args []string) error {
	command := exec.CommandContext(ctx, path, args...)
	command.Dir = options.Directory
	command.Stdin = options.Input
	command.Stdout = options.Output
	command.Stderr = options.Error

	return internal.Trace(ctx, "pyenv.run", func() error {
		return command.Run()
	}, "venv", e.venvDirectory, "path", path, "args", args)
}
