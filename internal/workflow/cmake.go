package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"fmt"
	"path/filepath"
)

// A cmakeTester is a special project tester that will invoke cmake before ctest
type cmakeTester struct {
	cmakeTool tool.CMake
	ctestTool tool.CTest
}

// Try to detect format tool for CMake project
func detectCMakeFormatTool(tools internal.Tools) internal.ProjectFormatter {
	if clangFormat := internal.GetTool[*tool.ClangFormat](tools); clangFormat.IsAvailable() {
		return clangFormat
	}

	return nil
}

// Try to detect lint tool for CMake project
func detectCMakeLintTool(tools internal.Tools) internal.ProjectLinter {
	if clangTidy := internal.GetTool[*tool.ClangTidy](tools); clangTidy.IsAvailable() {
		return clangTidy
	}

	return nil
}

// DetectCMakeProject detects if the project in the directory is a CMake project
func DetectCMakeProject(config internal.Config, tools internal.Tools) *internal.Project {
	if !internal.PathExists(filepath.Join(config.RootDirectory, "CMakeLists.txt")) {
		return nil
	}

	cmake := internal.GetTool[*tool.CMake](tools)

	if !cmake.IsAvailable() {
		return nil
	}

	formatTool := detectCMakeFormatTool(tools)
	lintTool := detectCMakeLintTool(tools)

	tester := &cmakeTester{
		cmakeTool: *cmake,
		ctestTool: *internal.GetTool[*tool.CTest](tools),
	}

	workflow := internal.Workflow{
		Configurator: cmake,
		Builder:      cmake,
		Tester:       tester,
		Formatter:    formatTool,
		Linter:       lintTool,
		Runner:       cmake,
	}

	var buildDirectory string

	if config.BuildDirectory != nil {
		buildDirectory = *config.BuildDirectory
	} else {
		buildDirectory = filepath.Join("build", "dev")
	}

	project := internal.MakeProject(config.RootDirectory, buildDirectory, cmake, workflow)

	return &project
}

func (t *cmakeTester) TestAll(project internal.Project, args []string) error {
	if err := t.cmakeTool.BuildAll(project, []string{}); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return t.ctestTool.Run(project, args)
}

func (t *cmakeTester) Test(project internal.Project, pattern string, args []string) error {
	if err := t.cmakeTool.BuildAll(project, []string{}); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return t.ctestTool.Run(project, append(args, "-R", pattern))
}
