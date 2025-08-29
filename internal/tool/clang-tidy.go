package tool

import (
	"cdt/internal"
	"fmt"
	"os"
	"path/filepath"
)

type ClangTidy struct {
	internal.ExecutableTool
}

// DetectClangTidy create a tool for clang-tidy
func DetectClangTidy(environment internal.Environment) *ClangTidy {
	return NewClangTidy(func() *internal.Executable {
		return environment.FindExecutable("clang-tidy")
	})
}

// NewClangTidy creates a clang-tidy tool from a custom executable
func NewClangTidy(detect func() *internal.Executable) *ClangTidy {
	return &ClangTidy{
		ExecutableTool: internal.MakeExecutableTool(
			"clang-tidy",
			"Clang Tidy",
			"A clang-based C++ “linter” tool.",
			internal.Tags{internal.ToolTagC, internal.ToolTagCpp, internal.ToolTagLint},
			detect,
		),
	}
}

func (c *ClangTidy) buildPaths(directory string, filenames []string) []string {
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

func (c *ClangTidy) buildArgs(rootDirectory string, buildDirectory string, paths []string) []string {
	var args []string

	configFile := filepath.Join(rootDirectory, ".clang-tidy")

	if _, err := os.Stat(configFile); err == nil {
		args = append(args, fmt.Sprintf("--config-file=%v", configFile))
	}

	args = append(args, "-p", buildDirectory)

	return append(args, paths...)
}

func (c *ClangTidy) LintAll(project internal.Project, args []string) error {
	info, err := project.Structure()

	if err != nil {
		return fmt.Errorf("failed to obtain project structure: %v", err)
	}

	paths := c.buildPaths(project.RootDirectory(), info.GetFiles())

	toolArgs := c.buildArgs(project.RootDirectory(), project.BuildDirectory(), paths)

	return c.ExecutableTool.RunForProject(project, append(toolArgs, args...))
}

func (c *ClangTidy) LintFiles(project internal.Project, filenames []string, args []string) error {
	paths := c.buildPaths(project.RootDirectory(), filenames)

	toolArgs := c.buildArgs(project.RootDirectory(), project.BuildDirectory(), paths)

	return c.ExecutableTool.RunForProject(project, append(toolArgs, args...))
}

func (c *ClangTidy) RunForProject(project internal.Project, args []string) error {
	toolArgs := []string{
		project.RootDirectory(),
	}

	return c.ExecutableTool.RunForProject(project, append(toolArgs, args...))
}
