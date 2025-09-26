package tool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"cdt/internal"
)

// A Go is a tool that wraps golang main tool `go`.
type Go struct {
	internal.ExecutableTool
}

// NewGo creates a go tool from a custom executable.
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

// DetectGo create go tool can be used in the project.
func DetectGo(ctx context.Context, environment internal.Environment) *Go {
	return NewGo(func() *internal.Executable {
		return environment.FindExecutable(ctx, "go")
	})
}

func (g *Go) Structure(
	ctx context.Context,
	info internal.ProjectInfo,
) (*internal.ProjectStructure, error) {
	structure := internal.ProjectStructure{
		Targets: make(map[string]internal.ProjectTarget),
	}

	builder := bytes.Buffer{}
	options := internal.RunOptions{Directory: info.Directory, Output: &builder, Error: nil}

	if err := g.Run(ctx, options, []string{"list", "-json=ImportPath,GoFiles", "./..."}); err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(&builder)
	for decoder.More() {
		var jsonData struct {
			ImportPath string   `json:"ImportPath"`
			GoFiles    []string `json:"GoFiles"`
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

func (g *Go) BuildAll(ctx context.Context, options internal.ProjectBuilderOptions) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "build"))
}

func (g *Go) BuildTargets(
	ctx context.Context,
	options internal.ProjectBuilderOptions,
	targets []string,
) error {
	return g.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append(options.ExtraArgs, "build"), targets...),
	)
}

func (g *Go) RunTarget(
	ctx context.Context,
	options internal.ProjectRunnerOptions,
	target string,
) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "run", target))
}

func (g *Go) TestAll(ctx context.Context, options internal.ProjectTesterOptions) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "test", "./..."))
}

func (g *Go) TestPattern(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "test", pattern))
}

func (g *Go) FormatAll(ctx context.Context, options internal.ProjectFormatterOptions) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "fmt", "./..."))
}

func (g *Go) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
	filenames []string,
) error {
	return g.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append(options.ExtraArgs, "fmt"), filenames...),
	)
}

func (g *Go) FormatCheckAll(_ context.Context, _ internal.ProjectFormatterOptions) error {
	return errors.New("go fmt doesn't support check mode")
}

func (g *Go) FormatCheckFiles(
	_ context.Context,
	_ internal.ProjectFormatterOptions,
	_ []string,
) error {
	return errors.New("go fmt doesn't support check mode")
}

func (g *Go) LintAll(ctx context.Context, options internal.ProjectLinterOptions) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "vet", "./..."))
}

func (g *Go) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	filenames []string,
) error {
	return g.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append(options.ExtraArgs, "vet"), filenames...),
	)
}

func (g *Go) AddDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	_ bool,
) error {
	return g.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append(options.ExtraArgs, "get"), dependencies...),
	)
}

func (g *Go) RemoveDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	_ bool,
) error {
	var noneDependencies []string

	for _, dependency := range dependencies {
		noneDependencies = append(noneDependencies, dependency+"@none")
	}

	return g.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append(options.ExtraArgs, "get"), noneDependencies...),
	)
}

func (g *Go) UpdateDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
) error {
	return g.RunForProject(
		ctx,
		options.ProjectInfo,
		append(append(options.ExtraArgs, "get"), dependencies...),
	)
}

func (g *Go) FetchDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
	_ bool,
) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "mod", "tidy"))
}

func (g *Go) ListDependencies(
	ctx context.Context,
	options internal.ProjectDependencyManagerOptions,
) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "list", "-m", "all"))
}

func (g *Go) AuditDependencies(
	_ context.Context,
	_ internal.ProjectDependencyManagerOptions,
) error {
	return errors.New("not supported")
}
