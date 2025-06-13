package tool

import (
	. "cdt/internal"
	"fmt"
	"os"
	"path/filepath"
)

const clangTidyMinVersion = 1
const clangTidyMaxVersion = 22

type ClangTidy struct {
	executable *Executable
	Version    *int
}

func detectClangTidyVersion(version int) *Executable {
	return FindExecutable(fmt.Sprintf("clang-tidy-%d", version))
}

// DetectClangTidy create a tool for clang-tidy
func DetectClangTidy(preferredVersion *int) *ClangTidy {
	if preferredVersion != nil {
		if executable := detectClangTidyVersion(*preferredVersion); executable != nil {
			return NewClangTidy(executable, preferredVersion)
		}
	}

	// Detect unversioned (system default)
	if executable := FindExecutable("clang-tidy"); executable != nil {
		return NewClangTidy(executable, nil)
	}

	for version := clangTidyMaxVersion; version >= clangTidyMinVersion; version-- {
		if executable := detectClangTidyVersion(version); executable != nil {
			return NewClangTidy(executable, &version)
		}
	}

	return NewClangTidy(nil, nil)
}

// NewClangTidy creates a clang-tidy tool from a custom executable
func NewClangTidy(executable *Executable, version *int) *ClangTidy {
	return &ClangTidy{
		executable: executable,
		Version:    version,
	}
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

// LintAll lints all files in the project
func (c *ClangTidy) LintAll(project Project, args []string) error {
	info, err := project.Structure()

	if err != nil {
		return err
	}

	paths := c.buildPaths(project.RootDirectory(), info.GetFiles())

	toolArgs := c.buildArgs(project.RootDirectory(), project.BuildDirectory(), paths)

	return c.executable.Run(append(toolArgs, args...))
}

// LintFiles lints a file in the project
func (c *ClangTidy) LintFiles(project Project, filenames []string, args []string) error {
	paths := c.buildPaths(project.RootDirectory(), filenames)

	toolArgs := c.buildArgs(project.RootDirectory(), project.BuildDirectory(), paths)

	return c.executable.Run(append(toolArgs, args...))
}

func (c *ClangTidy) Run(project Project, args []string) error {
	toolArgs := []string{
		project.RootDirectory(),
	}

	return c.executable.Run(append(toolArgs, args...))
}
