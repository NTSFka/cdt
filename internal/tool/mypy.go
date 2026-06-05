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
	options DetectOptions,
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
	args := options.ExtraArgs

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	} else {
		args = append(args, "*.py")
	}

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.LintReportFormatDefault:
		fallthrough
	case internal.LintReportFormatRaw:
		break
	default:
		return fmt.Errorf("unsupported report format: %s", options.Output.Format)
	}

	return m.RunForProject(ctx, options.ProjectInfo, args)
}
