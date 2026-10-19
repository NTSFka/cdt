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
	phpstan := internal.GetTool[*tool.PHPStan](tools)
	phpmd := internal.GetTool[*tool.PHPMD](tools)
	phpParallelLint := internal.GetTool[*tool.PHPParallelLint](tools)
	phpCsFixer := internal.GetTool[*tool.PHPCSFixer](tools)
	mago := internal.GetTool[*tool.Mago](tools)
	composerDependencyAnalyzer := internal.GetTool[*tool.ComposerDependencyAnalyzer](tools)
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
		Formatter:         &FormatterList{mago, phpCsFixer},
		Analyzer:          &AnalyzerList{phpstan, phpmd, mago, composerDependencyAnalyzer},
		Linter:            &LinterList{phpParallelLint, mago},
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
