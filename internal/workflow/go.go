package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"path/filepath"
)

// DetectGoProject detects if the project in the directory is a go project
func DetectGoProject(config internal.Config, tools tool.SupportedTools) *internal.Project {
	if !internal.PathExists(filepath.Join(config.RootDirectory, "go.mod")) {
		return nil
	}

	workflow := internal.Workflow{
		Configurator: nil,
		Builder:      tools.Go,
		Runner:       tools.Go,
		Tester:       tools.Go,
		Formatter:    tools.Go,
		Linter:       tools.GolangCILint,
	}

	project := internal.MakeProject(config.RootDirectory, "", tools.Go, workflow)

	return &project
}
