package project

import (
	. "cdt/internal"
	"cdt/internal/tool"
	"errors"
	"path/filepath"
)

// A cmakeTester is a special project tester that will invoke cmake before ctest
type cmakeTester struct {
	cmakeTool tool.CMake
	ctestTool tool.CTest
}

// Try to detect format tool for CMake project
func detectCMakeFormatTool(tools Tools) ProjectFormatter {
	if clangFormat := GetTool[*tool.ClangFormat](tools); clangFormat.IsAvailable() {
		return clangFormat
	}

	return nil
}

// Try to detect lint tool for CMake project
func detectCMakeLintTool(tools Tools) ProjectLinter {
	if clangTidy := GetTool[*tool.ClangTidy](tools); clangTidy.IsAvailable() {
		return clangTidy
	}

	return nil
}

// DetectCMakeProject detects if the project in the directory is a CMake project
func DetectCMakeProject(directory string, tools Tools) (ProjectStructureProvider, *Workflow, error) {
	if !PathExists(filepath.Join(directory, "CMakeLists.txt")) {
		return nil, nil, nil
	}

	cmake := GetTool[*tool.CMake](tools)

	if !cmake.IsAvailable() {
		return nil, nil, errors.New("cmake is not installed on the system")
	}

	formatTool := detectCMakeFormatTool(tools)
	lintTool := detectCMakeLintTool(tools)

	tester := &cmakeTester{
		cmakeTool: *cmake,
		ctestTool: *GetTool[*tool.CTest](tools),
	}

	return cmake, &Workflow{
		Configurator: cmake,
		Builder:      cmake,
		Tester:       tester,
		Formatter:    formatTool,
		Linter:       lintTool,
		Runner:       cmake,
	}, nil
}

func (t *cmakeTester) TestAll(project Project, args []string) error {
	if err := t.cmakeTool.BuildAll(project, []string{}); err != nil {
		return err
	}

	return t.ctestTool.Run(project, args)
}

func (t *cmakeTester) Test(project Project, pattern string, args []string) error {
	if err := t.cmakeTool.BuildAll(project, []string{}); err != nil {
		return err
	}

	return t.ctestTool.Run(project, append(args, "-R", pattern))
}
