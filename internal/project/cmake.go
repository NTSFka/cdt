package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"fmt"
	"path/filepath"
)

type CMakeType struct{}

func (c *CMakeType) Id() string {
	return "cmake"
}

func (c *CMakeType) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "CMakeLists.txt"))
}

func (c *CMakeType) Create(config Config, tools internal.Tools) internal.Project {
	cmake := internal.GetTool[*tool.CMake](tools)
	ctest := internal.GetTool[*tool.CTest](tools)
	clangFormat := internal.GetTool[*tool.ClangFormat](tools)
	clangTidy := internal.GetTool[*tool.ClangTidy](tools)

	tester := &cmakeTester{
		cmakeTool: cmake,
		ctestTool: ctest,
	}

	workflow := internal.Workflow{
		Configurator: cmake,
		Builder:      cmake,
		Tester:       tester,
		Formatter:    clangFormat,
		Linter:       clangTidy,
		Runner:       cmake,
	}

	var buildDirectory string

	if config.BuildDirectory != nil {
		buildDirectory = *config.BuildDirectory
	} else {
		buildDirectory = filepath.Join("build", "dev")
	}

	return internal.MakeProject(config.Directory, buildDirectory, cmake, workflow)
}

// A cmakeTester is a special project tester that will invoke cmake before ctest
type cmakeTester struct {
	cmakeTool *tool.CMake
	ctestTool *tool.CTest
}

func (t *cmakeTester) TestAll(project internal.Project, args []string) error {
	if err := t.cmakeTool.BuildAll(project, []string{}); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return t.ctestTool.RunForProject(project, args)
}

func (t *cmakeTester) Test(project internal.Project, pattern string, args []string) error {
	if err := t.cmakeTool.BuildAll(project, []string{}); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return t.ctestTool.RunForProject(project, append(args, "-R", pattern))
}
