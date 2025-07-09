package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"fmt"
	"path/filepath"
)

// A cmakeTester is a special project tester that will invoke cmake before ctest
type cmakeTester struct {
	cmakeTool *tool.CMake
	ctestTool *tool.CTest
}

// DetectCMakeProject detects if the project in the directory is a CMake project
func DetectCMakeProject(config internal.Config, tools tool.SupportedTools) *internal.Project {
	if !internal.PathExists(filepath.Join(config.RootDirectory, "CMakeLists.txt")) {
		return nil
	}

	tester := &cmakeTester{
		cmakeTool: tools.CMake,
		ctestTool: tools.CTest,
	}

	workflow := internal.Workflow{
		Configurator: tools.CMake,
		Builder:      tools.CMake,
		Tester:       tester,
		Formatter:    tools.ClangFormat,
		Linter:       tools.ClangTidy,
		Runner:       tools.CMake,
	}

	var buildDirectory string

	if config.BuildDirectory != nil {
		buildDirectory = *config.BuildDirectory
	} else {
		buildDirectory = filepath.Join("build", "dev")
	}

	project := internal.MakeProject(config.RootDirectory, buildDirectory, tools.CMake, workflow)

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
