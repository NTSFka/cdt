package tool

import (
	"cdt/internal"
	"context"
	"fmt"
)

const IdPyTest = "pytest"

type PyTest struct {
	internal.ExecutableTool
}

// DetectPyTest create a tool for pytest.
func DetectPyTest(
	ctx context.Context,
	options internal.DetectOptions,
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

	if options.Output.Reporter != nil {
		return fmt.Errorf("output to reporter not supported")
	}

	return p.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}
