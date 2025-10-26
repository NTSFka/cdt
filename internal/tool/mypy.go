package tool

import (
	"cdt/internal"
	"context"
)

const IdMyPy = "mypy"

type MyPy struct {
	internal.ExecutableTool
}

// DetectMyPy create a tool for mypy.
func DetectMyPy(
	ctx context.Context,
	config internal.ConfigTools,
	environment internal.Environment,
) *MyPy {
	return NewMyPy(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, config.Get(IdMyPy, "mypy"))
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

func (m *MyPy) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return m.RunForProject(ctx, options.ProjectInfo, append([]string{"*.py"}, options.ExtraArgs...))
}

func (m *MyPy) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	return m.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, filenames...))
}
