package internal

import (
	"context"
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

// NewExecutableTool creates an executable tool
func NewExecutableTool(id string, name string, info string, detect func() *Executable) *ExecutableTool {
	executable := MakeExecutableTool(id, name, info, detect)

	return &executable
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
		Assert(t.detect != nil, "detect function is not set")

		t.executable = t.detect()
		t.detected = true
	}

	return t.executable
}

func (t *ExecutableTool) RunContext(ctx context.Context, options RunOptions, args []string) error {
	if t.Executable() == nil {
		return t.NotFoundError()
	}

	return t.Executable().Run(ctx, options, args)
}

func (t *ExecutableTool) Run(project Project, args []string) error {
	return t.RunContext(context.Background(), RunOptions{Directory: project.RootDirectory()}, args)
}

// Tools is a container for available tools
type Tools []Tool

// Active returns only tools that are available
func (t *Tools) Active() (result Tools) {
	for _, tool := range *t {
		if tool.IsAvailable() {
			result = append(result, tool)
		}
	}

	return
}

// Get returns a tool by ID
func (t *Tools) Get(id string) Tool {
	for _, tool := range *t {
		if tool.Id() == id {
			return tool
		}
	}

	return nil
}

// PrintTable prints tools list to the writer
func (t *Tools) PrintTable(writer io.Writer) {
	if len(*t) == 0 {
		return
	}

	w := table.NewWriter()
	w.SetOutputMirror(writer)
	w.AppendHeader(table.Row{"ID", "Name", "Available", "Info", "Executable"})

	for _, tool := range *t {
		var executable string

		if path := tool.ExecutablePath(); path != nil {
			executable = *path
		} else {
			executable = "-"
		}

		w.AppendRow(table.Row{tool.Id(), tool.Name(), tool.IsAvailable(), tool.Info(), executable})
	}

	w.Render()
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
