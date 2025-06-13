package internal

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
)

// A Tool is an interface for Tools
type Tool interface {
	// The Id returns tool unique identifier
	Id() string

	// Name returns tool name
	Name() string

	// Info Get tool information
	Info() string

	// ExecutablePath returns tool executable path
	ExecutablePath() *string

	// IsAvailable return if the tool is available on the system
	IsAvailable() bool

	// Run Execute the tool on the project
	Run(project Project, args []string) error
}

// DependencyTool is an interface for a project dependencies management tool
type DependencyTool interface {
	Tool

	// InstallDependencies installs all defined dependencies
	InstallDependencies(project Project, args []string) error
}

// ConfigureTool is an interface for project configuration
type ConfigureTool interface {
	Tool

	// ConfigureProject configures the project
	ConfigureProject(project Project, args []string) error
}

// BuildTool is an interface for project builders
type BuildTool interface {
	Tool

	// BuildAll builds the whole project
	BuildAll(project Project, args []string) error

	// GetBuildTargets returns available targets to build in the project
	GetBuildTargets(project Project) ([]string, error)

	// BuildTarget builds a target from the project
	BuildTarget(project Project, target string, args []string) error
}

// TestTool is an interface for project testers
type TestTool interface {
	Tool

	// TestAll runs all tests in the project
	TestAll(project Project, args []string) error

	// Test runs tests that match the pattern
	Test(project Project, pattern string, args []string) error
}

// FormatTool is an interface for project formatters
type FormatTool interface {
	Tool

	// FormatAll formats all files in the project
	FormatAll(project Project, args []string) error

	// FormatFiles formats a file in the project
	FormatFiles(project Project, filenames []string, args []string) error

	// FormatCheckAll checks all files if some needs formatting
	FormatCheckAll(project Project, args []string) error

	// FormatCheckFiles checks a file if it needs formatting
	FormatCheckFiles(project Project, filenames []string, args []string) error
}

type Tools []Tool

func (t *Tools) Active() (result []Tool) {
	for _, tool := range *t {
		if tool.IsAvailable() {
			result = append(result, tool)
		}
	}

	return
}

// PrintToolList prints tools list to the writer
func PrintToolList(writer io.Writer, tools Tools) {
	if len(tools) == 0 {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(writer)
	t.AppendHeader(table.Row{"ID", "Name", "Available", "Info", "Executable"})

	for _, tool := range tools {
		var executable string

		if path := tool.ExecutablePath(); path != nil {
			executable = *path
		} else {
			executable = "-"
		}

		t.AppendRow(table.Row{tool.Id(), tool.Name(), tool.IsAvailable(), tool.Info(), executable})
	}

	t.Render()
}

// GetTool return a tool with required type
func GetTool[T Tool](tools Tools) T {
	for _, tool := range tools {
		if t, ok := tool.(T); ok {
			return t
		}
	}

	panic("Tool not found")
}
