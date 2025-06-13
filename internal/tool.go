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
