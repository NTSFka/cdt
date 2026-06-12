package tool

import (
	"context"
	"fmt"

	"cdt/internal"
)

const IdNilAway = "nilaway"

// A NilAway is a tool that wraps golang main tool `nilaway`.
type NilAway struct {
	internal.ExecutableTool
}

// NewNilAway creates a go tool from a custom executable.
func NewNilAway(detect internal.ExecutableToolDetectFunc) *NilAway {
	return &NilAway{
		internal.MakeExecutableTool(
			IdNilAway,
			"NilAway",
			"NilAway is a static analysis tool that seeks to help developers avoid nil panics in production.",
			internal.Tags{internal.ToolTagGo, internal.ToolTagLint},
			detect,
		),
	}
}

// DetectNilAway create nilaway tool can be used in the project.
func DetectNilAway(
	ctx context.Context,
	options DetectOptions,
) *NilAway {
	return NewNilAway(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdNilAway, "nilaway"))
	})
}

func (c *NilAway) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := appendFiles(options.ExtraArgs, options.Filenames, internal.StrPtr("./..."))

	if a, err := c.argsBuildLintOutputFormat(options.Output.Format); err == nil {
		args = append(args, a...)
	} else {
		return err
	}

	if options.Output.Filename != nil {
		return c.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return c.RunForProject(ctx, options.ProjectInfo, args)
}

func (c *NilAway) argsBuildLintOutputFormat(format internal.LintReportFormat) ([]string, error) {
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
