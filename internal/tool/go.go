package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/ctrf-io/go-ctrf-json-reporter/ctrf"
	"github.com/ctrf-io/go-ctrf-json-reporter/reporter"
	"golang.org/x/sync/errgroup"
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
	options DetectOptions,
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
	var pattern string
	if options.Pattern != nil {
		pattern = *options.Pattern
	} else {
		pattern = "./..."
	}

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.TestsReportFormatDefault:
		fallthrough
	case internal.TestsReportFormatRaw:
		break
	case internal.TestsReportFormatRawEvents:
		return g.testPatternJson(ctx, options, pattern)
	case internal.TestsReportFormatCtrf:
		return g.testPatternCtrf(ctx, options, pattern)
	default:
		return fmt.Errorf("unsupported report format: %s", options.Output.Format)
	}

	return g.RunForProject(ctx, options.ProjectInfo, append(options.ExtraArgs, "test", pattern))
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

	return g.RunForProject(ctx, options.ProjectInfo, args)
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

func (g *Go) testPatternJson(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	var writer io.Writer = os.Stdout

	if options.Output.Filename != nil {
		file, err := os.OpenFile(*options.Output.Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}

		writer = file

		defer func() {
			_ = file.Close()
		}()
	}

	runOptions := internal.RunOptions{
		Directory: options.Directory,
		Output:    writer,
	}

	return g.Run(
		ctx,
		runOptions,
		append(options.ExtraArgs, "test", "-json", pattern),
	)
}

func (g *Go) testPatternCtrfProcess(
	reader io.Reader,
	outputFilename *string,
) error {
	env := &ctrf.Environment{
		AppName:    "cdt",
		OSPlatform: runtime.GOOS,
	}

	report, err := reporter.ParseTestResults(reader, false, env)
	if err != nil {
		return fmt.Errorf("error parsing test results: %w", err)
	}

	if outputFilename != nil {
		err = report.WriteFile(*outputFilename)
	} else {
		err = report.Write(os.Stdout, true)
	}

	if err != nil {
		return fmt.Errorf("error writing the report to file: %w", err)
	}

	var buildFailed bool

	if report.Results.Extra != nil {
		extraMap, isMap := report.Results.Extra.(map[string]any)
		if !isMap {
			err = fmt.Errorf("expected a map, but got %T instead", report.Results.Extra)

			return fmt.Errorf("error extracting report results: %w", err)
		}

		if _, ok := extraMap["buildFail"]; ok {
			buildFailed = true
		}

		if _, ok := extraMap["FailedBuild"]; ok {
			buildFailed = true
		}
	}

	if report.Results.Summary.Failed > 0 {
		buildFailed = true
	}

	if buildFailed {
		return errors.New("build failed")
	}

	return nil
}

func (g *Go) testPatternCtrf(
	ctx context.Context,
	options internal.ProjectTesterOptions,
	pattern string,
) error {
	// TODO: use Pipe from created process
	reader, writer, err := os.Pipe()

	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}

	defer func() {
		_ = reader.Close()
	}()
	defer func() {
		_ = writer.Close()
	}()

	runOptions := internal.RunOptions{
		Directory: options.Directory,
		Input:     os.Stdin,
		Output:    writer,
		Error:     os.Stderr,
	}

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return g.testPatternCtrfProcess(reader, options.Output.Filename)
	})

	group.Go(
		func() error {
			err := g.Run(ctx, runOptions, append(options.ExtraArgs, "test", "-json", pattern))

			_ = writer.Close()

			return err
		},
	)

	return group.Wait()
}
