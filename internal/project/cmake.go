package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"errors"
	"fmt"
	"path/filepath"
)

// A cmakeTester is a special project tester that will invoke cmake before ctest
type cmakeTester struct {
	cmakeTool *tool.CMake
	ctestTool *tool.CTest
}

// DetectCMakeProject detects if the project in the directory is a CMake project
func DetectCMakeProject(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	if !internal.PathExists(filepath.Join(config.RootDirectory, "CMakeLists.txt")) {
		return nil, errors.New("cmake workflow requires CMakeLists.txt to be present in the project directory")
	}

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

	project := internal.MakeProject(config.RootDirectory, buildDirectory, cmake, workflow)

	return &project, nil
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
