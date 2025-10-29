package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"cdt/internal"
)

const IdClangTidy = "clang-tidy"

type ClangTidy struct {
	internal.ExecutableTool
}

// DetectClangTidy create a tool for clang-tidy.
func DetectClangTidy(
	ctx context.Context,
	options DetectOptions,
) *ClangTidy {
	return NewClangTidy(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(
			ctx,
			options.GetToolPath(IdClangTidy, "clang-tidy"),
		)
	})
}

// NewClangTidy creates a clang-tidy tool from a custom executable.
func NewClangTidy(detect internal.ExecutableToolDetectFunc) *ClangTidy {
	return &ClangTidy{
		ExecutableTool: internal.MakeExecutableTool(
			IdClangTidy,
			"Clang Tidy",
			"A clang-based C++ “linter” tool.",
			internal.Tags{internal.ToolTagC, internal.ToolTagCpp, internal.ToolTagLint},
			detect,
		),
	}
}

func (c *ClangTidy) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	structure, err := options.Structure(ctx)

	if err != nil {
		return fmt.Errorf("failed to obtain project structure: %w", err)
	}

	if options.OutputDirectory == nil {
		return internal.ErrNoOutputDirectory
	}

	toolArgs := c.buildArgs(options.Directory, *options.OutputDirectory, structure.GetFiles())

	return c.ExecutableTool.RunForProject(
		ctx,
		options.ProjectInfo,
		append(toolArgs, options.ExtraArgs...),
	)
}

func (c *ClangTidy) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	if options.OutputDirectory == nil {
		return internal.ErrNoOutputDirectory
	}

	toolArgs := c.buildArgs(options.Directory, *options.OutputDirectory, filenames)

	return c.ExecutableTool.RunForProject(
		ctx,
		options.ProjectInfo,
		append(toolArgs, options.ExtraArgs...),
	)
}

func (c *ClangTidy) RunForProject(
	ctx context.Context,
	info internal.ProjectInfo,
	args []string,
) error {
	toolArgs := []string{
		info.Directory,
	}

	return c.ExecutableTool.RunForProject(ctx, info, append(toolArgs, args...))
}

func (c *ClangTidy) buildArgs(
	rootDirectory string,
	buildDirectory string,
	paths []string,
) []string {
	var args []string

	configFile := filepath.Join(rootDirectory, ".clang-tidy")

	if _, err := os.Stat(configFile); err == nil {
		args = append(args, fmt.Sprintf("--config-file=%v", configFile))
	}

	args = append(args, "-p", buildDirectory)

	return append(args, paths...)
}
