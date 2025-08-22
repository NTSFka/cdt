package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"path/filepath"
)

type PythonType struct{}

func (p *PythonType) Id() string {
	return "python"
}

func (p *PythonType) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "pyproject.toml"))
}

func (p *PythonType) Create(config Config, tools internal.Tools) internal.Project {
	python := internal.GetTool[*tool.Python](tools)
	pyTest := internal.GetTool[*tool.PyTest](tools)
	pip := internal.GetTool[*tool.Pip](tools)

	workflow := internal.Workflow{
		Configurator:      nil,
		Builder:           nil,
		Runner:            python,
		Tester:            pyTest,
		Formatter:         nil,
		Linter:            nil,
		DependencyManager: pip,
	}

	return internal.MakeProject(config.Directory, "", &internal.EmptyProjectStructureProvider{}, workflow)
}
