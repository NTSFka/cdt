package project

import (
	. "cdt/internal"
	"cdt/internal/tool"
	"errors"
	"path/filepath"
)

// A CMakeProject describes a CMake project and operations on it
type CMakeProject struct {
	// CMakeTool represents the CMake tool
	CMakeTool tool.CMake
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

	c := &CMakeProject{
		CMakeTool: *cmake,
	}

	return cmake, &Workflow{
		Configurator: cmake,
		Builder:      cmake,
		Tester:       c,
		Formatter:    formatTool,
		Linter:       lintTool,
		Runner:       cmake,
	}, nil
}

// TestAll run project tester for all tests
func (p *CMakeProject) TestAll(project Project, args []string) error {
	if err := p.CMakeTool.BuildAll(project, []string{}); err != nil {
		return err
	}

	ctest := FindExecutable("ctest")

	if ctest == nil {
		return errors.New("ctest is not installed on the system")
	}

	return ctest.Run([]string{
		"--test-dir", project.BuildDirectory(),
	})
}

// Test run project tester with tests that matches the given pattern
func (p *CMakeProject) Test(project Project, pattern string, args []string) error {
	if err := p.CMakeTool.BuildAll(project, []string{}); err != nil {
		return err
	}

	ctest := FindExecutable("ctest")

	if ctest == nil {
		return errors.New("ctest is not installed on the system")
	}

	return ctest.Run([]string{
		"--test-dir", project.BuildDirectory(),
		"-R", pattern,
	})
}
