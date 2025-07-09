package internal

import (
	"fmt"
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

// ExecutableTool is a simple implementation of the Tool interface.
type ExecutableTool struct {
	id         string
	name       string
	info       string
	detected   bool
	detect     func() *Executable
	executable *Executable
}

// MakeExecutableTool creates an executable tool
func MakeExecutableTool(id string, name string, info string, detect func() *Executable) ExecutableTool {
	return ExecutableTool{
		id:         id,
		name:       name,
		info:       info,
		detected:   false,
		detect:     detect,
		executable: nil,
	}
}

func (t *ExecutableTool) NotFoundError() error {
	return fmt.Errorf("%v is not installed on the system", t.Name())
}

func (t *ExecutableTool) Id() string {
	return t.id
}

func (t *ExecutableTool) Name() string {
	return t.name
}

func (t *ExecutableTool) Info() string {
	return t.info
}

func (t *ExecutableTool) ExecutablePath() *string {
	if executable := t.Executable(); executable != nil {
		return &executable.Path
	}

	return nil
}

func (t *ExecutableTool) IsAvailable() bool {
	return t.Executable() != nil
}

func (t *ExecutableTool) Executable() *Executable {
	if !t.detected {
		if t.detect == nil {
			panic("detect function is not set")
		}

		t.executable = t.detect()
		t.detected = true
	}

	return t.executable
}

func (t *ExecutableTool) RunContext(ctx RunContext, args []string) error {
	if t.Executable() == nil {
		return t.NotFoundError()
	}

	return t.Executable().Run(ctx, args)
}

func (t *ExecutableTool) Run(project Project, args []string) error {
	return t.RunContext(NewRunContext(project.RootDirectory()), args)
}

// Tools is a container for available tools
type Tools []Tool

// Active returns only tools that are available
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

// GetTool return a tool with the required type
func GetTool[T Tool](tools Tools) T {
	for _, tool := range tools {
		if t, ok := tool.(T); ok {
			return t
		}
	}

	panic("Tool not found")
}
