package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"errors"
	"path/filepath"
)

// DetectGoProject detects if the project in the directory is a go project
func DetectGoProject(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	if !internal.PathExists(filepath.Join(config.RootDirectory, "go.mod")) {
		return nil, errors.New("go workflow requires go.mod to be present in the project directory")
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

	return &project, nil
}
