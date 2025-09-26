package tool

import (
	"cdt/internal"
	"context"
	"path/filepath"
)

type MyPy struct {
	internal.ExecutableTool
}

// DetectMyPy create a tool for mypy.
func DetectMyPy(ctx context.Context, environment internal.Environment) *MyPy {
	return NewMyPy(func() *internal.Executable {
		return environment.FindExecutable(ctx, "mypy")
	})
}

// NewMyPy creates a mypy tool from a custom executable.
func NewMyPy(detect func() *internal.Executable) *MyPy {
	return &MyPy{
		ExecutableTool: internal.MakeExecutableTool(
			"mypy",
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

func (m *MyPy) LintFiles(ctx context.Context, options internal.ProjectLinterOptions, filenames []string) error {
	paths := m.buildPaths(options.Directory, filenames)

	return m.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, paths...))
}

func (m *MyPy) buildPaths(directory string, filenames []string) []string {
	var paths []string

	for _, filename := range filenames {
		if filepath.IsAbs(filename) {
			paths = append(paths, filename)
		} else {
			paths = append(paths, filepath.Join(directory, filename))
		}
	}

	return paths
}
