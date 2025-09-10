package tool

import (
	"cdt/internal"
	"context"
)

type ParaTest struct {
	internal.ExecutableTool
}

// DetectParaTest create a tool for paratest
func DetectParaTest(ctx context.Context, environment internal.Environment) *ParaTest {
	return NewParaTest(func() *internal.Executable {
		// Detect composer vendor
		if executable := environment.FindExecutable(ctx, "vendor/bin/paratest"); executable != nil {
			return executable
		}

		// Detect unversioned (system default)
		if executable := environment.FindExecutable(ctx, "paratest"); executable != nil {
			return executable
		}

		return nil
	})
}

// NewParaTest creates a paratest tool from a custom executable
func NewParaTest(detect func() *internal.Executable) *ParaTest {
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

func (p *ParaTest) TestAll(ctx context.Context, info internal.ProjectInfo, args []string) error {
	return p.RunForProject(ctx, info, args)
}

func (p *ParaTest) Test(ctx context.Context, info internal.ProjectInfo, pattern string, args []string) error {
	return p.RunForProject(ctx, info, append(args, pattern))
}
