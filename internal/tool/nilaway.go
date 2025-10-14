package tool

import (
	"context"

	"cdt/internal"
)

// A NilAway is a tool that wraps golang main tool `nilaway`.
type NilAway struct {
	internal.ExecutableTool
}

// NewNilAway creates a go tool from a custom executable.
func NewNilAway(detect internal.ExecutableToolDetectFunc) *NilAway {
	return &NilAway{
		internal.MakeExecutableTool(
			"nilaway",
			"NilAway",
			"NilAway is a static analysis tool that seeks to help developers avoid nil panics in production.",
			internal.Tags{internal.ToolTagGo, internal.ToolTagLint},
			detect,
		),
	}
}

// DetectNilAway create nilaway tool can be used in the project.
func DetectNilAway(ctx context.Context, environment internal.Environment) *NilAway {
	return NewNilAway(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, "nilaway")
	})
}

func (c *NilAway) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return c.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "./..."))
}

func (c *NilAway) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	modules []string,
) error {
	return c.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, modules...))
}
