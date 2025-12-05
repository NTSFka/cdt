package workflow

import (
	"path/filepath"

	"cdt/internal"
	"cdt/internal/tool"
)

type PHP struct{}

func (p *PHP) Id() string {
	return "php"
}

func (p *PHP) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "composer.json"))
}

func (p *PHP) Create(config Config, tools internal.Tools) internal.Project {
	php := internal.GetTool[*tool.PHP](tools)
	phpStan := internal.GetTool[*tool.PHPStan](tools)
	phpCsFixer := internal.GetTool[*tool.PHPCSFixer](tools)
	paratest := internal.GetTool[*tool.ParaTest](tools)
	phpUnit := internal.GetTool[*tool.PHPUnit](tools)
	composer := internal.GetTool[*tool.Composer](tools)

	workflow := internal.Workflow{
		Name:              p.Id(),
		Configurator:      nil,
		Builder:           nil,
		Runner:            php,
		Tester:            &TesterFallback{paratest, phpUnit},
		CoverageCollector: &CoverageCollectorFallback{paratest, phpUnit},
		Formatter:         phpCsFixer,
		Linter:            phpStan,
		DependencyManager: composer,
	}

	return internal.Project{
		Info: internal.ProjectInfo{
			Directory:         config.Directory,
			StructureProvider: &internal.EmptyProjectStructureProvider{},
		},
		Workflow: workflow,
	}
}
