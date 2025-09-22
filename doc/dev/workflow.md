
# Workflow

Supported workflows are stored in the `workflow` package. The workflow type must implement the `workflow.Type` interface.
The interface requires `Id` function to return the unique identifier (not enforced) of the workflow, `Detect` function
should return if the given directory contains some files that matches the workflow, and `Create` function
should create a new workflow instance.

For building complex workflows, it's possible to use adapters like `workflow.LinterFallback` or `workflow.LinterList`.
A fallback adapter will run tools in the order until one is available. This is for cases when only one type of tool can
be executed (e.g., formatter). A list adapter will run all available tools in the list.

## Registration

The workflow type must be registered in the `workflow.SupportedTypes` variable (`workflow/type.go`).

## Example

```go
package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"path/filepath"
)

type Go struct {
}

func (g *Go) Id() string {
	return "go"
}

func (g *Go) Detect(directory string) bool {
	// If the project directory contains a go.mod file, it's a Go project
	return internal.PathExists(filepath.Join(directory, "go.mod"))
}

func (g *Go) Create(config Config, tools internal.Tools) internal.Project {
	goTool := internal.GetTool[*tool.Go](tools)
	goLint := internal.GetTool[*tool.GolangCILint](tools)

	workflow := internal.Workflow{
		Name:         g.Id(),
		Configurator: nil,
		Builder:      goTool,
		Runner:       goTool,
		Tester:       goTool,
		// Run the first available formatter
		Formatter:    &FormatterFallback{goTool, goLint},
		// Run all linters in the list if they are available
		Linter:       &LinterList{goTool, goLint},
	}

	return internal.Project{
		Info: internal.ProjectInfo{
			Directory:         config.Directory,
			StructureProvider: goTool,
		},
		Workflow: workflow,
	}
}
```