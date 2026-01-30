package workflow

import (
	"context"
	"fmt"
	"path/filepath"

	"cdt/internal"
	"cdt/internal/tool"
)

type CMake struct{}

func (c *CMake) Id() string {
	return "cmake"
}

func (c *CMake) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "CMakeLists.txt"))
}

func (c *CMake) Create(config Config, tools internal.Tools) internal.Project {
	cmake := internal.GetTool[*tool.CMake](tools)
	ctest := internal.GetTool[*tool.CTest](tools)
	clangFormat := internal.GetTool[*tool.ClangFormat](tools)
	clangTidy := internal.GetTool[*tool.ClangTidy](tools)

	tester := &cmakeTester{
		cmakeTool: cmake,
		ctestTool: ctest,
	}

	workflow := internal.Workflow{
		Name:         c.Id(),
		Configurator: cmake,
		Builder:      cmake,
		Tester:       tester,
		Formatter:    clangFormat,
		Linter:       clangTidy,
		Runner:       cmake,
	}

	var buildDirectory string

	if config.OutputDirectory != nil {
		buildDirectory = *config.OutputDirectory
	} else {
		buildDirectory = filepath.Join("build", "dev")
	}

	return internal.Project{
		Info: internal.ProjectInfo{
			Directory:         config.Directory,
			OutputDirectory:   &buildDirectory,
			StructureProvider: cmake,
		},
		Workflow: workflow,
	}
}

// A cmakeTester is a special project tester that will invoke cmake before ctest.
type cmakeTester struct {
	cmakeTool *tool.CMake
	ctestTool *tool.CTest
}

func (t *cmakeTester) Details() string {
	return "ctest"
}

func (t *cmakeTester) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return t.runTests(ctx, options, nil)
}

func (t *cmakeTester) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return t.runTests(ctx, options, &pattern)
}

func (t *cmakeTester) runTests(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern *string,
) error {
	if err := t.cmakeTool.BuildAll(ctx, internal.ProjectBuilderOptions{ProjectInfo: options.ProjectInfo}); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// TODO: other report formats
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
		args := options.ExtraArgs

		if pattern != nil {
			args = append(args, "-R", *pattern)
		}

		return t.ctestTool.RunForProject(ctx, options.ProjectInfo, args)
	case internal.TestsReportFormatEvents:
		break
	case internal.TestsReportFormatJson:
		break
	case internal.TestsReportFormatCtrf:
		break
	}

	return fmt.Errorf("unsupported report format: %s", options.Output.Format)
}
