package workflow

import (
	"path/filepath"

	"cdt/internal"
	"cdt/internal/tool"
)

type Go struct {
}

func (g *Go) Id() string {
	return "go"
}

func (g *Go) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "go.mod"))
}

func (g *Go) Create(config Config, tools internal.Tools) internal.Project {
	goTool := internal.GetTool[*tool.Go](tools)
	goLint := internal.GetTool[*tool.GolangCILint](tools)

	workflow := internal.Workflow{
		Name:              g.Id(),
		Configurator:      nil,
		Builder:           goTool,
		Runner:            goTool,
		Tester:            goTool,
		Formatter:         &FormatterFallback{goLint, goTool},
		Linter:            &LinterList{goTool, goLint},
		DependencyManager: goTool,
	}

	return internal.Project{
		Info: internal.ProjectInfo{
			Directory:         config.Directory,
			StructureProvider: goTool,
			OutputDirectory:   config.OutputDirectory,
		},
		Workflow: workflow,
	}
}
