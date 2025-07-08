package internal

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"os/exec"
)

// Environment represents an environment where the tools are located and can be executed.
type Environment interface {
	// The Id returns provider unique identifier
	Id() string

	// Start the environment
	Start() error

	// IsRunning returns if the environment is running
	IsRunning() bool

	// Stop the environment
	Stop() error

	// Cleanup cleans the environment
	Cleanup() error

	// FindExecutable try to find an executable in the environment
	FindExecutable(name string) *Executable

	// RunExecutable run an executable in the environment
	RunExecutable(ctx RunContext, path string, args []string) error
}

// EnvironmentProvider is an interface for providers that allows to create runtime environments.
type EnvironmentProvider interface {
	// The Id returns provider unique identifier
	Id() string

	// Name returns provider name
	Name() string

	// Info Get provider information
	Info() string

	// IsAvailable return if the provider is available on the system
	IsAvailable() bool

	// CreateEnvironment creates a new environment from parameter
	CreateEnvironment(directory string, parameter string) (Environment, error)
}

type EnvironmentProviders []EnvironmentProvider

// PrintList prints providers to the writer
func (p *EnvironmentProviders) PrintList(writer io.Writer) {
	if len(*p) == 0 {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.AppendHeader(table.Row{"ID", "Name", "Available", "Info"})

	for _, tool := range *p {
		t.AppendRow(table.Row{tool.Id(), tool.Name(), tool.IsAvailable(), tool.Info()})
	}

	t.Render()
}

// SystemEnvironment is the operating system environment
var SystemEnvironment Environment = &systemEnvironment{}

type systemEnvironment struct{}

func (s *systemEnvironment) Id() string {
	return "system"
}

func (s *systemEnvironment) Name() string {
	return "System environment"
}

func (s *systemEnvironment) Info() string {
	return "Default system environment"
}

func (s *systemEnvironment) IsAvailable() bool {
	// Always available
	return true
}

func (s *systemEnvironment) Start() error {
	// Always available
	return nil
}

func (s *systemEnvironment) IsRunning() bool {
	// Always available
	return true
}

func (s *systemEnvironment) Stop() error {
	// Always available
	return nil
}

func (s *systemEnvironment) Cleanup() error {
	// Always available
	return nil
}

func (s *systemEnvironment) FindExecutable(name string) *Executable {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil
	}

	return &Executable{Path: path, RunFunc: s.RunExecutable}
}

func (s *systemEnvironment) RunExecutable(ctx RunContext, path string, args []string) error {
	command := exec.Command(path, args...)
	command.Dir = ctx.Directory
	command.Stdout = ctx.Output
	command.Stderr = ctx.Error

	return command.Run()
}
