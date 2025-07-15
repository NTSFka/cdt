package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"path/filepath"
)

// DetectGoProject detects if the project in the directory is a go project
func DetectGoProject(config internal.Config, tools internal.Tools) *internal.Project {
	if !internal.PathExists(filepath.Join(config.RootDirectory, "go.mod")) {
		return nil
	}

	goTool := internal.GetTool[*tool.Go](tools)
	goLint := internal.GetTool[*tool.GolangCILint](tools)

	workflow := internal.Workflow{
		Configurator: nil,
		Builder:      goTool,
		Runner:       goTool,
		Tester:       goTool,
		Formatter:    goTool,
		Linter:       goLint,
	}

	project := internal.MakeProject(config.RootDirectory, "", goTool, workflow)

	return &project
}
