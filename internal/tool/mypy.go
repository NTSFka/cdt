package tool

import (
	"cdt/internal"
	"path/filepath"
)

type MyPy struct {
	internal.ExecutableTool
}

// DetectMyPy create a tool for mypy
func DetectMyPy(environment internal.Environment) *MyPy {
	return NewMyPy(func() *internal.Executable {
		return environment.FindExecutable("mypy")
	})
}

// NewMyPy creates a mypy tool from a custom executable
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

func (m *MyPy) LintAll(info internal.ProjectInfo, args []string) error {
	return m.RunForProject(info, append([]string{"*.py"}, args...))
}

func (m *MyPy) LintFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	paths := m.buildPaths(info.Directory, filenames)

	return m.RunForProject(info, append(args, paths...))
}
