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

func (p *ParaTest) RunTests(ctx context.Context, options internal.ProjectTesterOptions) error {
	// TODO: inline
	return p.runTests(ctx, options)
}

func (p *ParaTest) CollectCoverage(
	ctx context.Context,
	options internal.ProjectCoverageCollectorOptions,
) error {
	args := options.ExtraArgs

	filename := internal.DefaultString(options.Output.Filename, "php://stdout")

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.CoverageReportFormatDefault:
		fallthrough
	case internal.CoverageReportFormatRaw:
		args = append(args, "--coverage-text")
	case internal.CoverageReportFormatCobertura:
		args = append(args, "--coverage-cobertura", filename)
	case internal.CoverageReportFormatCrap4j:
		args = append(args, "--coverage-crap4j", filename)
	case internal.CoverageReportFormatHtml:
		args = append(args, "--coverage-html", filename)
	case internal.CoverageReportFormatXml:
		args = append(args, "--coverage-xml", filename)
	default:
		return fmt.Errorf("unsupported coverage report format: %s", options.Output.Format)
	}

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	return p.RunForProjectWithEnv(
		ctx,
		options.ProjectInfo,
		[]string{"XDEBUG_MODE=coverage"},
		args,
	)
}

func (p *ParaTest) runTests(
	ctx context.Context,
	options internal.ProjectTesterOptions,
) error {
	args := options.ExtraArgs

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	}

	filename := internal.DefaultString(options.Output.Filename, "php://stdout")

	// TODO: other report formats
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
		return p.RunForProject(ctx, options.ProjectInfo, args)
	case internal.TestsReportFormatRawEvents:
		return p.RunForProject(
			ctx,
			options.ProjectInfo,
			append(args, "--log-events-text", filename),
		)
	case internal.TestsReportFormatJson:
		break
	case internal.TestsReportFormatCtrf:
		break
	case internal.TestsReportFormatJUnit:
		return p.RunForProject(
			ctx,
			options.ProjectInfo,
			append(args, "--log-junit", filename),
		)
	case internal.TestsReportFormatTeamCity:
		return p.RunForProject(
			ctx,
			options.ProjectInfo,
			append(args, "--log-teamcity", filename),
		)
	}

	return fmt.Errorf("unsupported report format: %s", options.Output.Format)
}
