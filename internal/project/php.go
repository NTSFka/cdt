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

func (p *PHPType) Create(config Config, tools internal.Tools) internal.Project {
	php := internal.GetTool[*tool.PHP](tools)
	phpStan := internal.GetTool[*tool.PHPStan](tools)
	phpCsFixer := internal.GetTool[*tool.PHPCSFixer](tools)
	paratest := internal.GetTool[*tool.ParaTest](tools)
	phpUnit := internal.GetTool[*tool.PHPUnit](tools)

	workflow := internal.Workflow{
		Configurator: nil,
		Builder:      nil,
		Runner:       php,
		Tester:       &phpTester{paratest: paratest, phpunit: phpUnit},
		Formatter:    phpCsFixer,
		Linter:       phpStan,
	}

	return internal.MakeProject(config.Directory, "", &internal.EmptyProjectStructureProvider{}, workflow)
}

type phpTester struct {
	paratest *tool.ParaTest
	phpunit  *tool.PHPUnit
}

func (p *phpTester) TestAll(project internal.Project, args []string) error {
	if p.paratest.IsAvailable() {
		return p.paratest.TestAll(project, args)
	}

	return p.phpunit.TestAll(project, args)
}

func (p *phpTester) Test(project internal.Project, pattern string, args []string) error {
	if p.paratest.IsAvailable() {
		return p.paratest.Test(project, pattern, args)
	}

	return p.phpunit.Test(project, pattern, args)
}
