package tool

import (
	"context"
	"fmt"

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

func (p *PyTest) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return p.runTests(ctx, options, nil)
}

func (p *PyTest) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return p.runTests(ctx, options, &pattern)
}

func (p *PyTest) runTests(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern *string,
) error {
	// TODO: other report formats
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
		args := options.ExtraArgs

		if pattern != nil {
			args = append(args, *pattern)
		}

		return p.RunForProject(ctx, options.ProjectInfo, args)
	case internal.TestsReportFormatJson:
		break
	case internal.TestsReportFormatCtrf:
		break
	}

	return fmt.Errorf("unknown report format: %s", options.Output.Format)
}
