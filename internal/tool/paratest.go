package tool

import (
	"context"
	"fmt"
	"path/filepath"

	"cdt/internal"
)

const IdParaTest = "paratest"

type ParaTest struct {
	internal.ExecutableTool
}

// DetectParaTest create a tool for paratest.
func DetectParaTest(
	ctx context.Context,
	options DetectOptions,
) *ParaTest {
	if path, ok := options.ToolsPaths[IdParaTest]; ok {
		return NewParaTest(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewParaTest(internal.DetectExecutableChain(
		[]string{filepath.Join(options.ProjectDirectory, "vendor/bin/paratest"), "paratest"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewParaTest creates a paratest tool from a custom executable.
func NewParaTest(detect internal.ExecutableToolDetectFunc) *ParaTest {
	return &ParaTest{
		ExecutableTool: internal.MakeExecutableTool(
			IdParaTest,
			"ParaTest",
			"The objective of ParaTest is to support parallel testing in PHPUnit.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagTest},
			detect,
		),
	}
}

func (p *ParaTest) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return p.runTests(ctx, options, nil)
}

func (p *ParaTest) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return p.runTests(ctx, options, &pattern)
}

func (p *ParaTest) CollectCoverageAll(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	return p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		append(options.ExtraArgs, "--coverage-text"),
	)
}

func (p *ParaTest) CollectCoveragePattern(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
	pattern string,
) error {
	return p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		append(options.ExtraArgs, "--coverage-text", pattern),
	)
}

func (p *ParaTest) runTests(
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
