package tool

import (
	. "cdt/internal"
	"fmt"
	"os"
	"path/filepath"
)

type ClangTidy struct {
	executable *Executable
	Version    *int
}

func detectClangTidyVersion(version int) *Executable {
	return FindExecutable(fmt.Sprintf("clang-tidy-%d", version))
}

// NewClangTidy create a tool for clang-tidy
func NewClangTidy(preferredVersion *int) ClangTidy {
	if preferredVersion != nil {
		if executable := detectClangTidyVersion(*preferredVersion); executable != nil {
			return ClangTidy{executable: executable, Version: preferredVersion}
		}
	}

	// Detect unversioned (system default)
	if executable := FindExecutable("clang-tidy"); executable != nil {
		return ClangTidy{executable: executable, Version: nil}
	}

	// TODO: find a way to generate automatically
	versions := []int{22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	for _, version := range versions {
		if executable := detectClangTidyVersion(version); executable != nil {
			return ClangTidy{executable: executable, Version: &version}
		}
	}

	return ClangTidy{executable: nil, Version: nil}
}

func (c *ClangTidy) Id() string {
	return "clang-tidy"
}

func (c *ClangTidy) Name() string {
	return "Clang Tidy"
}

func (c *ClangTidy) Info() string {
	return "A clang-based C++ “linter” tool."
}

func (c *ClangTidy) ExecutablePath() *string {
	if c.executable != nil {
		return &c.executable.Path
	}

	return nil
}

func (c *ClangTidy) IsAvailable() bool {
	return c.executable != nil
}

func (c *ClangTidy) Enabled(directory string) bool {
	return PathExists(filepath.Join(directory, ".clang-tidy"))
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
	args = append(args, paths...)

	return args
}

// LintAll lints all files in the project
func (c *ClangTidy) LintAll(project Project, args []string) error {
	info, err := project.Structure()

	if err != nil {
		return err
	}

	paths := c.buildPaths(project.RootDirectory(), info.GetFiles())

	toolArgs := c.buildArgs(project.RootDirectory(), project.BuildDirectory(), paths)
	toolArgs = append(toolArgs, args...)

	return c.executable.Run(toolArgs)
}

// LintFiles lints a file in the project
func (c *ClangTidy) LintFiles(project Project, filenames []string, args []string) error {
	paths := c.buildPaths(project.RootDirectory(), filenames)

	toolArgs := c.buildArgs(project.RootDirectory(), project.BuildDirectory(), paths)
	toolArgs = append(toolArgs, args...)

	return c.executable.Run(toolArgs)
}

func (c *ClangTidy) Run(project Project, args []string) error {
	var toolArgs []string
	toolArgs = append(toolArgs, project.RootDirectory())
	toolArgs = append(toolArgs, args...)

	return c.executable.Run(toolArgs)
}
