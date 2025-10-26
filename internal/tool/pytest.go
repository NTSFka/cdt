package tool

import (
	"context"

	"cdt/internal"
)

const IdPyTest = "pytest"

type PyTest struct {
	internal.ExecutableTool
}

// DetectPyTest create a tool for pytest.
func DetectPyTest(
	ctx context.Context,
	config internal.ConfigTools,
	environment internal.Environment,
) *PyTest {
	return NewPyTest(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, config.Get(IdPyTest, "pytest"))
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
	return p.RunForProject(ctx, options.ProjectInfo, options.ExtraArgs)
}

func (p *PyTest) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, pattern))
}
