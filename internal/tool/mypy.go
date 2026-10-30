package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdMyPy = "mypy"

type MyPy struct {
	internal.ExecutableTool
}

// DetectMyPy create a tool for mypy.
func DetectMyPy(
	ctx context.Context,
	options internal.DetectOptions,
) *MyPy {
	return NewMyPy(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdMyPy, "mypy"))
	})
}

// NewMyPy creates a mypy tool from a custom executable.
func NewMyPy(detect internal.ExecutableToolDetectFunc) *MyPy {
	return &MyPy{
		ExecutableTool: internal.MakeExecutableTool(
			IdMyPy,
			"MyPy",
			"Mypy is a static type checker for Python.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (m *MyPy) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := appendFiles(options.ExtraArgs, options.Filenames, new("*.py"))

	if options.Output.Reporter != nil {
		return fmt.Errorf("output to reporter not supported")
	}

	return m.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}
