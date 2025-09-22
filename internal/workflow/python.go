package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"path/filepath"
)

type Python struct{}

func (p *Python) Id() string {
	return "python"
}

func (p *Python) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "pyproject.toml"))
}

func (p *Python) Create(config Config, tools internal.Tools) internal.Project {
	python := internal.GetTool[*tool.Python](tools)
	pyTest := internal.GetTool[*tool.PyTest](tools)
	pip := internal.GetTool[*tool.Pip](tools)
	pylint := internal.GetTool[*tool.Pylint](tools)
	flake := internal.GetTool[*tool.Flake8](tools)
	mypy := internal.GetTool[*tool.MyPy](tools)
	ruff := internal.GetTool[*tool.Ruff](tools)
	bandit := internal.GetTool[*tool.Bandit](tools)
	black := internal.GetTool[*tool.Black](tools)

	workflow := internal.Workflow{
		Name:              p.Id(),
		Configurator:      nil,
		Builder:           nil,
		Runner:            python,
		Tester:            pyTest,
		Formatter:         &FormatterFallback{ruff, black},
		Linter:            &LinterList{pylint, flake, mypy, ruff, bandit},
		DependencyManager: pip,
	}

	return internal.Project{
		Info: internal.ProjectInfo{
			Directory:         config.Directory,
			StructureProvider: &internal.EmptyProjectStructureProvider{},
		},
		Workflow: workflow,
	}
}
