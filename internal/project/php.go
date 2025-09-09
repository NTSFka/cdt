package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"path/filepath"
)

type PHPType struct{}

func (p *PHPType) Id() string {
	return "php"
}

func (p *PHPType) Detect(directory string) bool {
	return internal.PathExists(filepath.Join(directory, "composer.json"))
}

func (p *PHPType) Create(config Config, tools internal.Tools) Project {
	php := internal.GetTool[*tool.PHP](tools)
	phpStan := internal.GetTool[*tool.PHPStan](tools)
	phpCsFixer := internal.GetTool[*tool.PHPCSFixer](tools)
	paratest := internal.GetTool[*tool.ParaTest](tools)
	phpUnit := internal.GetTool[*tool.PHPUnit](tools)
	composer := internal.GetTool[*tool.Composer](tools)

	workflow := internal.Workflow{
		Configurator:      nil,
		Builder:           nil,
		Runner:            php,
		Tester:            &TesterFallback{paratest, phpUnit},
		Formatter:         phpCsFixer,
		Linter:            phpStan,
		DependencyManager: composer,
	}

	return Project{
		Desc:     internal.Project{Directory: config.Directory, StructureProvider: &internal.EmptyProjectStructureProvider{}},
		Workflow: workflow,
	}
}
