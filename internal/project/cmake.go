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
	// A RootDirectory is a path to the project directory
	rootDirectory string

	// A buildDirectory where the project is being built.
	buildDirectory string

	// CMakeTool represents the CMake tool
	CMakeTool tool.CMake

	FormatTool FormatTool
	LintTool   LintTool
}

// Try to detect format tool for CMake project
func detectCMakeFormatTool(tools Tools) FormatTool {
	if clangFormat := GetTool[*tool.ClangFormat](tools); clangFormat.IsAvailable() {
		return clangFormat
	}

	return nil
}

// Try to detect lint tool for CMake project
func detectCMakeLintTool(tools Tools) LintTool {
	if clangTidy := GetTool[*tool.ClangTidy](tools); clangTidy.IsAvailable() {
		return clangTidy
	}

	return nil
}

// DetectCMakeProject detects if the project in the directory is a CMake project
func DetectCMakeProject(directory string, buildDirectory string, tools Tools) (*CMakeProject, error) {
	if !PathExists(filepath.Join(directory, "CMakeLists.txt")) {
		return nil, nil
	}

	cmake := GetTool[*tool.CMake](tools)

	if !cmake.IsAvailable() {
		return nil, errors.New("cmake is not installed on the system")
	}

	formatTool := detectCMakeFormatTool(tools)
	lintTool := detectCMakeLintTool(tools)

	return &CMakeProject{
		rootDirectory:  directory,
		buildDirectory: buildDirectory,
		CMakeTool:      *cmake,
		FormatTool:     formatTool,
		LintTool:       lintTool,
	}, nil
}

func (p *CMakeProject) RootDirectory() string {
	return p.rootDirectory
}

func (p *CMakeProject) BuildDirectory() string {
	return p.buildDirectory
}

func (p *CMakeProject) Structure(project Project) (*ProjectStructure, error) {
	if err := p.Configure(project); err != nil {
		return nil, err
	}

	fileApi := utils.NewCmakeFileApi(p.buildDirectory)

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
func (p *CMakeProject) Configure(project Project) error {
	fileApi := utils.NewCmakeFileApi(p.buildDirectory)

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
func (p *CMakeProject) BuildAll(project Project) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	return p.CMakeTool.BuildAll(project, []string{})
}

// BuildTarget builds a specific target
func (p *CMakeProject) BuildTarget(project Project, target string) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	return p.CMakeTool.BuildTargets(project, []string{target}, []string{})
}

// BuildTargets builds specific targets
func (p *CMakeProject) BuildTargets(project Project, targets []string) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	return p.CMakeTool.BuildTargets(project, targets, []string{})
}

// FormatAll formats all files in the CMake project
func (p *CMakeProject) FormatAll(project Project) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	if p.FormatTool != nil {
		return p.FormatTool.FormatAll(project, []string{})
	}

	return errors.New("no format tool found")
}

// FormatFiles formats a file in the project
func (p *CMakeProject) FormatFiles(project Project, filenames []string) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	if p.FormatTool != nil {
		return p.FormatTool.FormatFiles(project, filenames, []string{})
	}

	return errors.New("no format tool found")
}

// FormatCheckAll check all files in the project if some needs formatting
func (p *CMakeProject) FormatCheckAll(project Project) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	if p.FormatTool != nil {
		return p.FormatTool.FormatCheckAll(project, []string{})
	}

	return errors.New("no format tool found")
}

// FormatCheckFiles check a file in the project if it needs formatting
func (p *CMakeProject) FormatCheckFiles(project Project, filenames []string) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	if p.FormatTool != nil {
		return p.FormatTool.FormatCheckFiles(project, filenames, []string{})
	}

	return errors.New("no format tool found")
}

// LintAll lint all files in the project
func (p *CMakeProject) LintAll(project Project) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	if p.LintTool != nil {
		return p.LintTool.LintAll(project, []string{})
	}

	return nil
}

// LintFiles lint a file(s) in the project
func (p *CMakeProject) LintFiles(project Project, filenames []string) error {
	if err := p.Configure(project); err != nil {
		return err
	}

	if p.LintTool != nil {
		return p.LintTool.LintFiles(project, filenames, []string{})
	}

	return nil
}

// TestAll run project tester for all tests
func (p *CMakeProject) TestAll(project Project) error {
	if err := p.BuildAll(project); err != nil {
		return err
	}

	ctest := FindExecutable("ctest")

	if ctest == nil {
		return errors.New("ctest is not installed on the system")
	}

	return ctest.Run([]string{
		"--test-dir", p.buildDirectory,
	})
}

// Test run project tester with tests that matches the given pattern
func (p *CMakeProject) Test(project Project, pattern string) error {
	if err := p.BuildAll(project); err != nil {
		return err
	}

	ctest := FindExecutable("ctest")

	if ctest == nil {
		return errors.New("ctest is not installed on the system")
	}

	return ctest.Run([]string{
		"--test-dir", p.buildDirectory,
		"-R", pattern,
	})
}

func (p *CMakeProject) Run(project Project, target string, args []string) error {
	if err := p.BuildAll(project); err != nil {
		return err
	}

	fileApi := utils.NewCmakeFileApi(p.buildDirectory)

	reply, err := fileApi.Reply()
	if err != nil {
		return err
	}

	for _, t := range reply.Targets {
		if t.Name == target && t.Type == utils.TargetExecutable {
			executable := Executable{Path: filepath.Join(p.buildDirectory, t.Name)}

			return executable.Run(args)
		}
	}

	return fmt.Errorf("target '%s' not found", target)
}
