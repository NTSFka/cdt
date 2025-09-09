package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// A Go is a tool that wraps golang main tool `go`.
type Go struct {
	internal.ExecutableTool
}

// NewGo creates a go tool from a custom executable
func NewGo(detect func() *internal.Executable) *Go {
	return &Go{
		internal.MakeExecutableTool(
			"go",
			"Go",
			"tool for managing Go source code",
			internal.Tags{
				internal.ToolTagGo,
				internal.ToolTagBuild,
				internal.ToolTagRun,
				internal.ToolTagTest,
				internal.ToolTagFormat,
				internal.ToolTagLint,
			},
			detect,
		),
	}
}

// DetectGo create go tool can be used in the project
func DetectGo(environment internal.Environment) *Go {
	return NewGo(func() *internal.Executable {
		return environment.FindExecutable("go")
	})
}

func (g *Go) Structure(project internal.Project) (*internal.ProjectStructure, error) {
	info := internal.ProjectStructure{
		Targets: make(map[string]internal.ProjectTarget),
	}

	builder := bytes.Buffer{}
	options := internal.RunOptions{Directory: project.RootDirectory(), Output: &builder, Error: nil}
	if err := g.Run(context.Background(), options, []string{"list", "-json=ImportPath,GoFiles", "./..."}); err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(&builder)
	for decoder.More() {
		var jsonData struct {
			ImportPath string
			GoFiles    []string
		}
		if err := decoder.Decode(&jsonData); err != nil {
			return nil, fmt.Errorf("json decode failed: %w", err)
		}

		info.Targets[jsonData.ImportPath] = internal.ProjectTarget{
			Files: jsonData.GoFiles,
		}
	}

	return &info, nil
}

func (g *Go) BuildAll(project internal.Project, args []string) error {
	return g.RunForProject(project, append(args, "build"))
}

func (g *Go) BuildTargets(project internal.Project, targets []string, args []string) error {
	return g.RunForProject(project, append(append(args, "build"), targets...))
}

func (g *Go) RunTarget(project internal.Project, target string, args []string) error {
	return g.RunForProject(project, append(args, "run", target))
}

func (g *Go) TestAll(project internal.Project, args []string) error {
	return g.RunForProject(project, append(args, "test", "./..."))
}

func (g *Go) Test(project internal.Project, pattern string, args []string) error {
	return g.RunForProject(project, append(args, "test", pattern))
}

func (g *Go) FormatAll(project internal.Project, args []string) error {
	return g.RunForProject(project, append(args, "fmt", "./..."))
}

func (g *Go) FormatFiles(project internal.Project, filenames []string, args []string) error {
	return g.RunForProject(project, append(append(args, "fmt"), filenames...))
}

func (g *Go) FormatCheckAll(_ internal.Project, _ []string) error {
	return errors.New("go fmt doesn't support check mode")
}

func (g *Go) FormatCheckFiles(_ internal.Project, _ []string, _ []string) error {
	return errors.New("go fmt doesn't support check mode")
}

func (g *Go) LintAll(project internal.Project, args []string) error {
	return g.RunForProject(project, append(args, "vet", "./..."))
}

func (g *Go) LintFiles(project internal.Project, filenames []string, args []string) error {
	return g.RunForProject(project, append(append(args, "vet"), filenames...))
}
