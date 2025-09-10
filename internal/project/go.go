package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"path/filepath"
)

type GoType struct {
}

func (g *GoType) Id() string {
	return "go"
}

func (g *GoType) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "go.mod"))
}

func (g *GoType) Create(config Config, tools internal.Tools) internal.Project {
	goTool := internal.GetTool[*tool.Go](tools)
	goLint := internal.GetTool[*tool.GolangCILint](tools)

	workflow := internal.Workflow{
		Configurator: nil,
		Builder:      goTool,
		Runner:       goTool,
		Tester:       goTool,
		Formatter:    goTool,
		Linter:       &LinterList{goTool, goLint},
	}

	return internal.Project{
		Type:     "go",
		Desc:     internal.ProjectInfo{Directory: config.Directory, StructureProvider: goTool},
		Workflow: workflow,
	}
}
