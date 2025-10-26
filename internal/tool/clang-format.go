package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cdt/internal"
)

const IdClangFormat = "clang-format"

// ClangFormat is a formatter for the clang-format tool.
type ClangFormat struct {
	internal.ExecutableTool
}

// DetectClangFormat CreateEnvironment clang-format tool can be used in the project.
func DetectClangFormat(
	ctx context.Context,
	config internal.ConfigTools,
	environment internal.Environment,
) *ClangFormat {
	return NewClangFormat(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, config.Get(IdClangFormat, "clang-format"))
	})
}

// NewClangFormat creates a clang-format tool from a custom executable.
func NewClangFormat(detect internal.ExecutableToolDetectFunc) *ClangFormat {
	return &ClangFormat{
		ExecutableTool: internal.MakeExecutableTool(
			IdClangFormat,
			"Clang Format",
			"A tool to format C/C++/Java/JavaScript/JSON/Objective-C/Protobuf/C# code.",
			internal.Tags{internal.ToolTagC, internal.ToolTagCpp, internal.ToolTagFormat},
			detect,
		),
	}
}

// FormatAll formats all files in the project.
func (c *ClangFormat) FormatAll(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	structure, err := options.Structure(ctx)

	if err != nil {
		return fmt.Errorf("failed to obtain project structure: %w", err)
	}

	toolArgs := c.buildArgs(options.Directory, []string{"-i"}, structure.GetFiles())

	return c.RunForProject(ctx, options.ProjectInfo, append(toolArgs, options.ExtraArgs...))
}

// FormatFiles formats a file in the project.
func (c *ClangFormat) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	toolArgs := c.buildArgs(options.Directory, []string{"-i"}, filenames)

	return c.RunForProject(ctx, options.ProjectInfo, append(toolArgs, options.ExtraArgs...))
}

// FormatCheckAll checks all files if some needs formatting.
func (c *ClangFormat) FormatCheckAll(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	structure, err := options.Structure(ctx)

	if err != nil {
		return fmt.Errorf("failed to obtain project structure: %w", err)
	}

	toolArgs := c.buildArgs(options.Directory, []string{"--dry-run"}, structure.GetFiles())

	return c.RunForProject(ctx, options.ProjectInfo, append(toolArgs, options.ExtraArgs...))
}

// FormatCheckFiles checks a file if it needs formatting.
func (c *ClangFormat) FormatCheckFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	toolArgs := c.buildArgs(options.Directory, []string{"--dry-run"}, filenames)

	return c.RunForProject(ctx, options.ProjectInfo, append(toolArgs, options.ExtraArgs...))
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
