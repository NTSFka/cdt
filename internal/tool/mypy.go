package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdMyPy = "mypy"

type MyPy struct {
	internal.ExecutableTool
}

// DetectMyPy create a tool for mypy.
func DetectMyPy(
	ctx context.Context,
	options internal.DetectOptions,
) *MyPy {
	return NewMyPy(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdMyPy, "mypy"))
	})
}

// NewMyPy creates a mypy tool from a custom executable.
func NewMyPy(detect internal.ExecutableToolDetectFunc) *MyPy {
	return &MyPy{
		ExecutableTool: internal.MakeExecutableTool(
			IdMyPy,
			"MyPy",
			"Mypy is a static type checker for Python.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagLint},
			detect,
		),
	}
}

func (m *MyPy) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := appendFiles(options.ExtraArgs, options.Filenames, internal.StrPtr("*.py"))

	if a, err := m.argsBuildLintOutputFormat(options.Output.Format); err == nil {
		args = append(args, a...)
	} else {
		return err
	}

	if options.Output.Filename != nil {
		return m.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return m.RunForProject(ctx, options.ProjectInfo, args)
}

func (m *MyPy) argsBuildLintOutputFormat(format internal.LintReportFormat) ([]string, error) {
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
