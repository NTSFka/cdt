package project

import (
	. "cdt/internal"
	"cdt/internal/tool"
	"cdt/internal/utils"
	"errors"
	"fmt"
	"path/filepath"
)

// A CMakeProject describes a CMake project and operations on it
type CMakeProject struct {
	// CMakeTool represents the CMake tool
	CMakeTool tool.CMake
}

// Try to detect format tool for CMake project
func detectCMakeFormatTool(tools Tools) FormatTool {
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

	return c, &Workflow{
		Configurator: c,
		Builder:      c,
		Tester:       c,
		Formatter:    formatTool,
		Linter:       lintTool,
		Runner:       c,
	}, nil
}

func (p *CMakeProject) Structure(project Project) (*ProjectStructure, error) {
	if err := p.Configure(project, []string{}); err != nil {
		return nil, err
	}

	fileApi := utils.NewCmakeFileApi(project.BuildDirectory())

	info := ProjectStructure{
		Targets: make(map[string]ProjectTarget),
	}

	if reply, err := fileApi.Reply(); err == nil {
		for _, target := range reply.Targets {
			info.Targets[target.Name] = ProjectTarget{
				Files:      target.Files,
				Dependency: target.External,
			}
		}
	}

	return &info, nil
}

// Configure configures the CMake project
func (p *CMakeProject) Configure(project Project, args []string) error {
	fileApi := utils.NewCmakeFileApi(project.BuildDirectory())

	if err := fileApi.Query("codemodel", 2); err != nil {
		return err
	}

	res := p.CMakeTool.ConfigureProject(project, []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
	})

	if res != nil {
		return res
	}

	return res
}

// BuildAll builds all targets in the CMake project
func (p *CMakeProject) BuildAll(project Project, args []string) error {
	if err := p.Configure(project, []string{}); err != nil {
		return err
	}

	return p.CMakeTool.BuildAll(project, []string{})
}

// BuildTarget builds a specific target
func (p *CMakeProject) BuildTarget(project Project, target string, args []string) error {
	if err := p.Configure(project, []string{}); err != nil {
		return err
	}

	return p.CMakeTool.BuildTargets(project, []string{target}, []string{})
}

// BuildTargets builds specific targets
func (p *CMakeProject) BuildTargets(project Project, targets []string, args []string) error {
	if err := p.Configure(project, []string{}); err != nil {
		return err
	}

	return p.CMakeTool.BuildTargets(project, targets, []string{})
}

// TestAll run project tester for all tests
func (p *CMakeProject) TestAll(project Project, args []string) error {
	if err := p.BuildAll(project, []string{}); err != nil {
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
	if err := p.BuildAll(project, []string{}); err != nil {
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

func (p *CMakeProject) Run(project Project, target string, args []string) error {
	if err := p.BuildAll(project, []string{}); err != nil {
		return err
	}

	fileApi := utils.NewCmakeFileApi(project.BuildDirectory())

	reply, err := fileApi.Reply()
	if err != nil {
		return err
	}

	for _, t := range reply.Targets {
		if t.Name == target && t.Type == utils.TargetExecutable {
			executable := Executable{Path: filepath.Join(project.BuildDirectory(), t.Name)}

			return executable.Run(args)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}
