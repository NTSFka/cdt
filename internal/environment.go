package internal

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Environment represents an environment where the tools are located and can be executed.
type Environment interface {
	// The Id returns environment unique identifier
	Id() string

	// Start the environment
	Start(ctx context.Context) error

	// IsRunning returns if the environment is running
	IsRunning(ctx context.Context) bool

	// Stop the environment
	Stop(ctx context.Context) error

	// Cleanup cleans the environment
	Cleanup(ctx context.Context) error

	// FindExecutable try to find an executable in the environment
	FindExecutable(name string) *Executable

	// RunExecutable run an executable in the environment
	RunExecutable(ctx context.Context, options RunOptions, path string, args []string) error
}

// EnvironmentProvider is an interface for providers that allows to create runtime environments.
type EnvironmentProvider interface {
	// The Id returns provider unique identifier
	Id() string

	// Aliases returns provider alternative unique identifiers
	Aliases() []string

	// Detect try to detect an environment in the directory.
	Detect(directory string) *Environment

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

// PrintTable prints providers to the writer
func (p *EnvironmentProviders) PrintTable(writer io.Writer) {
	if len(*p) == 0 {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.AppendHeader(table.Row{"ID (aliases)", "Name", "Available", "Info"})

	for _, tool := range *p {
		var id string

		if len(tool.Aliases()) > 0 {
			id = fmt.Sprintf("%v (%v)", tool.Id(), strings.Join(tool.Aliases(), ", "))
		} else {
			id = tool.Id()
		}

		t.AppendRow(table.Row{
			id,
			tool.Name(),
			tool.IsAvailable(),
			tool.Info(),
		})
	}

	t.Render()
}

var SystemEnvironmentProvider EnvironmentProvider = &systemEnvironmentProvider{}

type systemEnvironmentProvider struct{}

func (s *systemEnvironmentProvider) Id() string {
	return "system"
}

func (s *systemEnvironmentProvider) Aliases() []string {
	return []string{"s"}
}

func (s *systemEnvironmentProvider) Detect(_ string) *Environment {
	// System environment is a special case that is always available, and there is no need to detect it
	return nil
}

func (s *systemEnvironmentProvider) Name() string {
	return "System"
}

func (s *systemEnvironmentProvider) Info() string {
	return "Native OS system environment"
}

func (s *systemEnvironmentProvider) IsAvailable() bool {
	// Always available
	return true
}

func (s *systemEnvironmentProvider) CreateEnvironment(_ string, _ string) (Environment, error) {
	return SystemEnvironment, nil
}

// SystemEnvironment is the operating system environment
var SystemEnvironment Environment = &systemEnvironment{}

type systemEnvironment struct{}

func (s *systemEnvironment) Id() string {
	return "system"
}

func (s *systemEnvironment) Start(_ context.Context) error {
	// Always available
	return nil
}

func (s *systemEnvironment) IsRunning(_ context.Context) bool {
	// Always available
	return true
}

func (s *systemEnvironment) Stop(_ context.Context) error {
	// Always available
	return nil
}

func (s *systemEnvironment) Cleanup(_ context.Context) error {
	// Always available
	return nil
}

func (s *systemEnvironment) FindExecutable(name string) *Executable {
	return Trace(context.Background(), "system.find_executable", func() *Executable {
		path, err := exec.LookPath(name)
		if err != nil {
			return nil
		}

		return &Executable{Path: path, Runtime: s}
	}, "name", name)
}

func (s *systemEnvironment) RunExecutable(ctx context.Context, options RunOptions, path string, args []string) error {
	command := exec.CommandContext(ctx, path, args...)
	command.Dir = options.Directory
	command.Stdin = options.Input
	command.Stdout = options.Output
	command.Stderr = options.Error

	return Trace(ctx, "system.run", func() error {
		return command.Run()
	}, "path", path, "args", args)
}
