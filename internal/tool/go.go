package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"runtime"

	"github.com/ctrf-io/go-ctrf-json-reporter/ctrf"
	"github.com/ctrf-io/go-ctrf-json-reporter/reporter"
)

const IdGo = "go"

// A Go is a tool that wraps golang main tool `go`.
type Go struct {
	internal.ExecutableTool
}

// NewGo creates a go tool from a custom executable.
func NewGo(detect internal.ExecutableToolDetectFunc) *Go {
	return &Go{
		internal.MakeExecutableTool(
			IdGo,
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
func DetectGo(
	ctx context.Context,
	options internal.DetectOptions,
) *Go {
	return NewGo(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdGo, "go"))
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

	if err := g.Run(
		ctx,
		options,
		[]string{"list", "-json=ImportPath,GoFiles", "./..."},
	); err != nil {
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

func (g *Go) BuildTargets(
	ctx context.Context,
	options internal.ProjectBuilderOptions,
) error {
	args := append([]string{"build"}, options.ExtraArgs...)

	if options.OutputDirectory != nil {
		args = append(args, "-o", *options.OutputDirectory)
	}

	if options.Targets != nil && len(*options.Targets) > 0 {
		args = append(args, *options.Targets...)
	}

	return g.RunForProject(ctx, options.ProjectInfo, args)
}

func (g *Go) RunTarget(
	ctx context.Context,
	options internal.ProjectRunnerOptions,
	target string,
) error {
	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "run", target))
}

func (g *Go) RunTests(ctx context.Context, options internal.ProjectTesterOptions) error {
	args := append([]string{"test"}, options.ExtraArgs...)

	if options.Pattern != nil {
		args = append(args, *options.Pattern)
	} else {
		args = append(args, "./...")
	}

	if options.Output.Reporter != nil {
		return g.runTestsReporter(ctx, options, options.Output.Reporter, args)
	}

	return g.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (g *Go) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	if options.CheckOnly {
		return errors.New("go fmt doesn't support check mode")
	}

	args := append([]string{"fmt"}, options.ExtraArgs...)

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	} else {
		args = append(args, "./...")
	}

	return g.RunForProject(
		ctx,
		options.ProjectInfo,
		args,
	)
}

func (g *Go) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"vet"}, options.ExtraArgs...)

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	} else {
		args = append(args, "./...")
	}

	if options.Output.Reporter != nil {
		return fmt.Errorf("output to reporter not supported")
	}

	return g.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
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

func (g *Go) runTestsReporter(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	reporter internal.TestsReporter,
	args []string,
) error {
	// Process JSON output
	args = append(args, "-json")

	outputBuffer := bytes.Buffer{}

	err := g.RunForProjectWithWriter(ctx, options.ProjectInfo, &outputBuffer, args)

	if err != nil && !checkExitCode(err, []int{1}) {
		return err
	} else if outputBuffer.Len() == 0 {
		// Nothing returned
		return nil
	}

	return errors.Join(err, g.reportJsonOutput(&outputBuffer, reporter))
}

func (g *Go) reportJsonOutput(
	reader io.Reader,
	testsReporter internal.TestsReporter,
) error {
	env := &ctrf.Environment{
		AppName:    "go",
		OSPlatform: runtime.GOOS,
	}

	report, err := reporter.ParseTestResults(reader, false, env)
	if err != nil {
		return fmt.Errorf("error parsing test results: %w", err)
	}

	return testsReporter.AddCtrf(report)
}
