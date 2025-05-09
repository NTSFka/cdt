package tool

import (
	. "cdt/internal"
	"fmt"
	"os"
	"path/filepath"
)

// ClangFormat is a formatter for the clang-format tool
type ClangFormat struct {
	executable *Executable
	Version    *int
}

func detectClangFormat() *Executable {
	return FindExecutable("clang-format")
}

func detectClangFormatVersion(version int) *Executable {
	return FindExecutable(fmt.Sprintf("clang-format-%d", version))
}

// NewClangFormat Create clang-format tool can be used in the project
func NewClangFormat(preferredVersion *int) ClangFormat {
	if preferredVersion != nil {
		if executable := detectClangFormatVersion(*preferredVersion); executable != nil {
			return ClangFormat{executable: executable, Version: preferredVersion}
		}
	}

	// Detect unversioned (system default)
	if executable := detectClangFormat(); executable != nil {
		return ClangFormat{executable: executable, Version: nil}
	}

	// TODO: find a way to generate automatically
	versions := []int{22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	for _, version := range versions {
		if executable := detectClangFormatVersion(version); executable != nil {
			return ClangFormat{executable: executable, Version: &version}
		}
	}

	return ClangFormat{executable: nil, Version: nil}
}

func (c *ClangFormat) Id() string {
	return "clang-format"
}

func (c *ClangFormat) Name() string {
	return "Clang Format"
}

func (c *ClangFormat) Info() string {
	return "A tool to format C/C++/Java/JavaScript/JSON/Objective-C/Protobuf/C# code."
}

func (c *ClangFormat) ExecutablePath() *string {
	if c.executable != nil {
		return &c.executable.Path
	}

	return nil
}

func (c *ClangFormat) IsAvailable() bool {
	return c.executable != nil
}

func (c *ClangFormat) Enabled(directory string) bool {
	return PathExists(filepath.Join(directory, ".clang-format"))
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

func (c *ClangFormat) buildArgs(directory string, paths []string) []string {
	var args []string

	configFile := filepath.Join(directory, ".clang-format")

	if _, err := os.Stat(configFile); err == nil {
		args = append(args, fmt.Sprintf("--style=file:%v", configFile))
	}

	args = append(args, "--Werror")
	args = append(args, paths...)

	return args
}

// FormatAll formats all files in the project
func (c *ClangFormat) FormatAll(project Project, args []string) error {
	info, err := project.Structure()

	if err != nil {
		return err
	}

	paths := c.buildPaths(project.RootDirectory(), info.GetFiles())

	toolArgs := c.buildArgs(project.RootDirectory(), paths)
	toolArgs = append(toolArgs, "-i")
	toolArgs = append(toolArgs, args...)

	return c.executable.Run(toolArgs)
}

// FormatFiles formats a file in the project
func (c *ClangFormat) FormatFiles(project Project, filenames []string, args []string) error {
	paths := c.buildPaths(project.RootDirectory(), filenames)

	toolArgs := c.buildArgs(project.RootDirectory(), paths)
	toolArgs = append(toolArgs, "-i")
	toolArgs = append(toolArgs, args...)

	return c.executable.Run(toolArgs)
}

// FormatCheckAll checks all files if some needs formatting
func (c *ClangFormat) FormatCheckAll(project Project, args []string) error {
	info, err := project.Structure()

	if err != nil {
		return err
	}

	paths := c.buildPaths(project.RootDirectory(), info.GetFiles())

	toolArgs := c.buildArgs(project.RootDirectory(), paths)
	toolArgs = append(toolArgs, "--dry-run")
	toolArgs = append(toolArgs, args...)

	return c.executable.Run(toolArgs)
}

// FormatCheckFiles checks a file if it needs formatting
func (c *ClangFormat) FormatCheckFiles(project Project, filenames []string, args []string) error {
	paths := c.buildPaths(project.RootDirectory(), filenames)

	toolArgs := c.buildArgs(project.RootDirectory(), paths)
	toolArgs = append(toolArgs, "--dry-run")
	toolArgs = append(toolArgs, args...)

	return c.executable.Run(toolArgs)
}

func (c *ClangFormat) Run(_ Project, args []string) error {
	return c.executable.Run(args)
}
