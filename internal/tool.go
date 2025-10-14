package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const (
	// Commands.
	ToolTagConfigure   Tag = "configure"
	ToolTagBuild       Tag = "build"
	ToolTagTest        Tag = "test"
	ToolTagLint        Tag = "lint"
	ToolTagFormat      Tag = "format"
	ToolTagRun         Tag = "run"
	ToolTagDependency  Tag = "dependency"
	ToolTagEnvironment Tag = "environment"
	// Languages.
	ToolTagGo     Tag = "go"
	ToolTagC      Tag = "c"
	ToolTagCpp    Tag = "c++"
	ToolTagPython Tag = "python"
	ToolTagPhp    Tag = "php"
)

type Tag string
type Tags []Tag

// A Tool is an interface for Tools.
type Tool interface {
	// The Id returns tool unique identifier
	Id() string

	// Name returns tool name
	Name() string

	// Info Get tool information
	Info() string

	// Tags returns tool tags that can be used to filter tools
	Tags() Tags

	// IsAvailable return if the tool is available on the system
	IsAvailable() bool

	// Executable returns tool executable
	Executable() *Executable

	// Run Execute the tool on the project
	Run(ctx context.Context, options RunOptions, args []string) error
}

type ExecutableToolDetectFunc func() (*Executable, error)

// ExecutableTool is a simple implementation of the Tool interface.
type ExecutableTool struct {
	id          string
	name        string
	info        string
	tags        Tags
	detected    bool
	detectError error
	detect      ExecutableToolDetectFunc
	executable  *Executable
}

// MakeExecutableTool creates an executable tool.
func MakeExecutableTool(
	id string,
	name string,
	info string,
	tags Tags,
	detect ExecutableToolDetectFunc,
) ExecutableTool {
	return ExecutableTool{
		id:          id,
		name:        name,
		info:        info,
		tags:        tags,
		detected:    false,
		detectError: nil,
		detect:      detect,
		executable:  nil,
	}
}

// NewExecutableTool creates an executable tool.
func NewExecutableTool(
	id string,
	name string,
	info string,
	tags Tags,
	detect ExecutableToolDetectFunc,
) *ExecutableTool {
	executable := MakeExecutableTool(id, name, info, tags, detect)

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

func (t *ExecutableTool) Tags() Tags {
	return t.tags
}

func (t *ExecutableTool) IsAvailable() bool {
	return t.Executable() != nil
}

func (t *ExecutableTool) Executable() *Executable {
	if !t.detected {
		Assert(t.detect != nil, "detect function is not set")

		t.executable, t.detectError = t.detect()
		t.detected = true
	}

	return t.executable
}

func (t *ExecutableTool) Run(ctx context.Context, options RunOptions, args []string) error {
	if t.Executable() == nil {
		if t.detectError != nil {
			return t.detectError
		}

		return t.NotFoundError()
	}

	return t.Executable().Run(ctx, options, args)
}

func (t *ExecutableTool) RunForProject(ctx context.Context, info ProjectInfo, args []string) error {
	options := RunOptions{
		Directory: info.Directory,
		Input:     os.Stdin,
		Output:    os.Stdout,
		Error:     os.Stderr,
	}

	return t.Run(ctx, options, args)
}

// Tools is a container for available tools.
type Tools []Tool

// OnlyAvailable returns only tools that are available.
func (t *Tools) OnlyAvailable() (result Tools) {
	if t == nil {
		return
	}

	for _, tool := range *t {
		if tool.IsAvailable() {
			result = append(result, tool)
		}
	}

	return
}

func (t *Tools) FilterByTags(tags []string) (result Tools) {
	Assert(len(tags) > 0, "tags is empty")

	if t == nil {
		return
	}

	for _, tool := range *t {
		contains := true

		for _, tag := range tags {
			contains = contains && slices.Contains(tool.Tags(), Tag(tag))
		}

		if contains {
			result = append(result, tool)
		}
	}

	return
}

// Get returns a tool by ID.
func (t *Tools) Get(id string) Tool {
	for _, tool := range *t {
		if tool.Id() == id {
			return tool
		}
	}

	return nil
}

// PrintTable prints tools list to the writer.
func (t *Tools) PrintTable(writer io.Writer) {
	if t == nil || len(*t) == 0 {
		return
	}

	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(writer)
	tableWriter.AppendHeader(table.Row{"ID", "Name", "Available", "Tags", "Executable", "Info"})

	if width := DetectTermWidth(writer); width != nil {
		tableWriter.SetAllowedRowLength(*width)
	}

	for _, tool := range *t {
		var executable string

		if exe := tool.Executable(); exe != nil {
			executable = fmt.Sprintf("%v:%v", exe.Runtime.Id(), exe.Path)
		} else {
			executable = "-"
		}

		var tags []string

		for _, tag := range tool.Tags() {
			tags = append(tags, string(tag))
		}

		tableWriter.AppendRow(table.Row{
			tool.Id(),
			tool.Name(),
			tool.IsAvailable(),
			strings.Join(tags, ", "),
			executable,
			tool.Info(),
		})
	}

	tableWriter.Render()
}

// GetTool return a tool with the required type.
func GetTool[T Tool](tools Tools) T {
	for _, tool := range tools {
		if t, ok := tool.(T); ok {
			return t
		}
	}

	panic("Tool not found")
}
