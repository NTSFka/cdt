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

func (c *Go) Structure(project internal.Project) (*internal.ProjectStructure, error) {
	info := internal.ProjectStructure{
		Targets: make(map[string]internal.ProjectTarget),
	}

	builder := bytes.Buffer{}
	options := internal.RunOptions{Directory: project.RootDirectory(), Output: &builder, Error: nil}
	if err := c.Run(context.Background(), options, []string{"list", "-json=ImportPath,GoFiles", "./..."}); err != nil {
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

func (c *Go) BuildAll(project internal.Project, args []string) error {
	return c.RunForProject(project, append(args, "build"))
}

func (c *Go) BuildTargets(project internal.Project, targets []string, args []string) error {
	return c.RunForProject(project, append(append(args, "build"), targets...))
}

func (c *Go) RunTarget(project internal.Project, target string, args []string) error {
	return c.RunForProject(project, append(args, "run", target))
}

func (c *Go) TestAll(project internal.Project, args []string) error {
	return c.RunForProject(project, append(args, "test", "./..."))
}

func (c *Go) Test(project internal.Project, pattern string, args []string) error {
	return c.RunForProject(project, append(args, "test", pattern))
}

func (c *Go) FormatAll(project internal.Project, args []string) error {
	return c.RunForProject(project, append(args, "fmt", "./..."))
}

func (c *Go) FormatFiles(project internal.Project, filenames []string, args []string) error {
	return c.RunForProject(project, append(append(args, "fmt"), filenames...))
}

func (c *Go) FormatCheckAll(_ internal.Project, _ []string) error {
	return errors.New("go fmt doesn't support check mode")
}

func (c *Go) FormatCheckFiles(_ internal.Project, _ []string, _ []string) error {
	return errors.New("go fmt doesn't support check mode")
}
