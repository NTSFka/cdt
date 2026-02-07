package tool

import (
	"context"
	"fmt"
	"runtime"

	"cdt/internal"
)

const IdPyTest = "pytest"

type PyTest struct {
	internal.ExecutableTool
}

// DetectPyTest create a tool for pytest.
func DetectPyTest(
	ctx context.Context,
	options DetectOptions,
) *PyTest {
	return NewPyTest(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdPyTest, "pytest"))
	})
}

// NewPyTest creates a pytest tool from a custom executable.
func NewPyTest(detect internal.ExecutableToolDetectFunc) *PyTest {
	return &PyTest{
		ExecutableTool: internal.MakeExecutableTool(
			IdPyTest,
			"PyTest",
			"The pytest framework makes it easy to write small, readable tests, and can scale to support complex "+
				"functional testing for applications and libraries.",
			internal.Tags{internal.ToolTagPython, internal.ToolTagTest},
			detect,
		),
	}
}

func (p *PyTest) RunTests(ctx context.Context, options internal.ProjectTesterOptions) error {
	args := options.ExtraArgs

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	filename := internal.DefaultString(options.Output.Filename, "/dev/stdout")

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
		break
	case internal.TestsReportFormatJUnit:
		if runtime.GOOS == "windows" && options.Output.Filename == nil {
			return fmt.Errorf("output to stdout is not supported on Windows")
		}

		args = append(args, "--junitxml="+filename)
	default:
		return fmt.Errorf("unsupported report format: %s", options.Output.Format)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
