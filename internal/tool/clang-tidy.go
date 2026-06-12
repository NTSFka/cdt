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

func (c *ClangTidy) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	if options.OutputDirectory == nil {
		return internal.ErrNoOutputDirectory
	}

	var filenames []string
	if options.Filenames != nil && len(*options.Filenames) > 0 {
		filenames = *options.Filenames
	} else {
		structure, err := options.Structure(ctx)
		if err != nil {
			return fmt.Errorf("failed to obtain project structure: %w", err)
		}

		filenames = structure.GetFiles()
	}

	args := append(
		c.buildArgs(options.Directory, *options.OutputDirectory, filenames),
		options.ExtraArgs...,
	)

	if a, err := c.argsBuildLintOutputFormat(options.Output.Format); err == nil {
		args = append(args, a...)
	} else {
		return err
	}

	if options.Output.Filename != nil {
		return c.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return c.ExecutableTool.RunForProject(ctx, options.ProjectInfo, args)
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

func (c *ClangTidy) argsBuildLintOutputFormat(format internal.LintReportFormat) ([]string, error) {
	var args []string

	// nolint: exhaustive
	switch format {
	case internal.LintReportFormatDefault:
		fallthrough
	case internal.LintReportFormatRaw:
		break
	default:
		return nil, fmt.Errorf("unsupported report format: %s", format)
	}

	return args, nil
}
