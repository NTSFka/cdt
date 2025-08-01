package tool

import (
	"cdt/internal"
	"fmt"
	"os"
	"path/filepath"
)

// ClangFormat is a formatter for the clang-format tool
type ClangFormat struct {
	internal.ExecutableTool
}

// DetectClangFormat CreateEnvironment clang-format tool can be used in the project
func DetectClangFormat(environment internal.Environment) *ClangFormat {
	return NewClangFormat(func() *internal.Executable {
		return environment.FindExecutable("clang-format")
	})
}

// NewClangFormat creates a clang-format tool from a custom executable
func NewClangFormat(detect func() *internal.Executable) *ClangFormat {
	return &ClangFormat{
		ExecutableTool: internal.MakeExecutableTool(
			"clang-format",
			"Clang Format",
			"A tool to format C/C++/Java/JavaScript/JSON/Objective-C/Protobuf/C# code.",
			detect,
		),
	}
}

func (c *ClangFormat) buildPaths(directory string, filenames []string) []string {
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

func (c *ClangFormat) buildArgs(directory string, extraArgs []string, paths []string) []string {
	var args []string

	configFile := filepath.Join(directory, ".clang-format")

	if _, err := os.Stat(configFile); err == nil {
		args = append(args, fmt.Sprintf("--style=file:%v", configFile))
	}

	args = append(args, "--Werror")
	args = append(args, extraArgs...)
	args = append(args, paths...)

	return args
}

// FormatAll formats all files in the project
func (c *ClangFormat) FormatAll(project internal.Project, args []string) error {
	info, err := project.Structure()

	if err != nil {
		return fmt.Errorf("failed to obtain project structure: %v", err)
	}

	paths := c.buildPaths(project.RootDirectory(), info.GetFiles())

	toolArgs := c.buildArgs(project.RootDirectory(), []string{"-i"}, paths)

	return c.RunForProject(project, append(toolArgs, args...))
}

// FormatFiles formats a file in the project
func (c *ClangFormat) FormatFiles(project internal.Project, filenames []string, args []string) error {
	paths := c.buildPaths(project.RootDirectory(), filenames)

	toolArgs := c.buildArgs(project.RootDirectory(), []string{"-i"}, paths)

	return c.RunForProject(project, append(toolArgs, args...))
}

// FormatCheckAll checks all files if some needs formatting
func (c *ClangFormat) FormatCheckAll(project internal.Project, args []string) error {
	info, err := project.Structure()

	if err != nil {
		return fmt.Errorf("failed to obtain project structure: %v", err)
	}

	paths := c.buildPaths(project.RootDirectory(), info.GetFiles())

	toolArgs := c.buildArgs(project.RootDirectory(), []string{"--dry-run"}, paths)

	return c.RunForProject(project, append(toolArgs, args...))
}

// FormatCheckFiles checks a file if it needs formatting
func (c *ClangFormat) FormatCheckFiles(project internal.Project, filenames []string, args []string) error {
	paths := c.buildPaths(project.RootDirectory(), filenames)

	toolArgs := c.buildArgs(project.RootDirectory(), []string{"--dry-run"}, paths)

	return c.RunForProject(project, append(toolArgs, args...))
}
