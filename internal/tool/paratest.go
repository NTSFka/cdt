package tool

import (
	"context"

	"cdt/internal"
)

type ParaTest struct {
	internal.ExecutableTool
}

// DetectParaTest create a tool for paratest.
func DetectParaTest(ctx context.Context, environment internal.Environment) *ParaTest {
	return NewParaTest(internal.DetectExecutableChain(
		[]string{"vendor/bin/paratest", "paratest"},
		func(name string) (*internal.Executable, error) {
			return environment.FindExecutable(ctx, name)
		},
	))
}

// NewParaTest creates a paratest tool from a custom executable.
func NewParaTest(detect internal.ExecutableToolDetectFunc) *ParaTest {
	return &ParaTest{
		ExecutableTool: internal.MakeExecutableTool(
			"paratest",
			"ParaTest",
			"The objective of ParaTest is to support parallel testing in PHPUnit.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagTest},
			detect,
		),
	}
}

func (p *ParaTest) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return p.RunForProject(ctx, options.ProjectInfo, options.ExtraArgs)
}

func (p *ParaTest) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return p.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, pattern))
}
