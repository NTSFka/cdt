package tool

import (
	"cdt/internal"
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type ClangTidy struct {
	internal.ExecutableTool
}

// DetectClangTidy create a tool for clang-tidy
func DetectClangTidy(ctx context.Context, environment internal.Environment) *ClangTidy {
	return NewClangTidy(func() *internal.Executable {
		return environment.FindExecutable(ctx, "clang-tidy")
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

func (c *ClangTidy) LintAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	structure, err := info.Structure(ctx)

	if err != nil {
		return fmt.Errorf("failed to obtain project structure: %v", err)
	}

	if info.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	paths := c.buildPaths(info.Directory, structure.GetFiles())

	toolArgs := c.buildArgs(info.Directory, *info.IntermediateDirectory, paths)

	return c.ExecutableTool.RunForProject(ctx, info, append(toolArgs, args...))
}

func (c *ClangTidy) LintFiles(ctx context.Context, info internal.ProjectInfo, filenames []string, args []string) error {
	paths := c.buildPaths(info.Directory, filenames)

	if info.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	toolArgs := c.buildArgs(info.Directory, *info.IntermediateDirectory, paths)

	return c.ExecutableTool.RunForProject(ctx, info, append(toolArgs, args...))
}

func (c *ClangTidy) RunForProject(ctx context.Context, info internal.ProjectInfo, args []string) error {
	toolArgs := []string{
		info.Directory,
	}

	return c.ExecutableTool.RunForProject(ctx, info, append(toolArgs, args...))
}
