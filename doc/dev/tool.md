
# Tool

The tools are the main functionality of the cdt tool. The minimal tool must implement the `Tool` interface. The easiest
way to implement this interface is to embed `internal.ExecutableTool`. The implementation of the tool must be stored in
`internal/tool` directory.

## Registration

The tool must be registered in the `InitTools` function in the `tool/tools.go` file. Because the tools are created for
a specific environment, a global instance cannot be created.

If tool is an environment provider, it must be registered in the `InitEnvironmentProviders` function.

## Tool functionality

The additional functionality can be added by implementing a specific interface defined by a workflow. There are workflow
interfaces like `ProjectLinter` or `ProjectFormatter`. If the tool implements a specific interface, it can be used in the
workflow or invoked directly in the command line. This adds a unified interface for invoking the tool for a specific task
like linting or formatting.

## Example

```go
package tool

import (
	"cdt/internal"
	"context"
)

type MyTool struct {
	internal.ExecutableTool
}

func DetectMyTool(ctx context.Context, environment internal.Environment) *MyTool {
	// Define lazy detectable function in the environment
	// This allows to detect the tool only when it is necessary
	return NewMyTool(func() *internal.Executable {
		return environment.FindExecutable(ctx, "my-tool")
	})
}

func NewMyTool(detect func() *internal.Executable) *MyTool {
	return &MyTool{
		ExecutableTool: internal.MakeExecutableTool(
			"my-tool",
			"MyTool",
			"My tool description",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

// Tool support linting

func (t *MyTool) LintFiles(ctx context.Context, options internal.ProjectLinterOptions) error {
	return t.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, *options.Filenames...))
}
```