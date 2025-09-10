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

func (g *Go) Structure(info internal.ProjectInfo) (*internal.ProjectStructure, error) {
	structure := internal.ProjectStructure{
		Targets: make(map[string]internal.ProjectTarget),
	}

	builder := bytes.Buffer{}
	options := internal.RunOptions{Directory: info.Directory, Output: &builder, Error: nil}
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

		structure.Targets[jsonData.ImportPath] = internal.ProjectTarget{
			Files: jsonData.GoFiles,
		}
	}

	return &structure, nil
}

func (g *Go) BuildAll(info internal.ProjectInfo, args []string) error {
	return g.RunForProject(info, append(args, "build"))
}

func (g *Go) BuildTargets(info internal.ProjectInfo, targets []string, args []string) error {
	return g.RunForProject(info, append(append(args, "build"), targets...))
}

func (g *Go) RunTarget(info internal.ProjectInfo, target string, args []string) error {
	return g.RunForProject(info, append(args, "run", target))
}

func (g *Go) TestAll(info internal.ProjectInfo, args []string) error {
	return g.RunForProject(info, append(args, "test", "./..."))
}

func (g *Go) Test(info internal.ProjectInfo, pattern string, args []string) error {
	return g.RunForProject(info, append(args, "test", pattern))
}

func (g *Go) FormatAll(info internal.ProjectInfo, args []string) error {
	return g.RunForProject(info, append(args, "fmt", "./..."))
}

func (g *Go) FormatFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return g.RunForProject(info, append(append(args, "fmt"), filenames...))
}

func (g *Go) FormatCheckAll(_ internal.ProjectInfo, _ []string) error {
	return errors.New("go fmt doesn't support check mode")
}

func (g *Go) FormatCheckFiles(_ internal.ProjectInfo, _ []string, _ []string) error {
	return errors.New("go fmt doesn't support check mode")
}

func (g *Go) LintAll(info internal.ProjectInfo, args []string) error {
	return g.RunForProject(info, append(args, "vet", "./..."))
}

func (g *Go) LintFiles(info internal.ProjectInfo, filenames []string, args []string) error {
	return g.RunForProject(info, append(append(args, "vet"), filenames...))
}
