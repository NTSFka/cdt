package tool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"cdt/internal"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

const IdGolangCILint = "golangci-lint"

// A GolangCILint is a tool that wraps golang main tool `golangci-lint`.
type GolangCILint struct {
	internal.ExecutableTool
}

// NewGolangCILint creates a go tool from a custom executable.
func NewGolangCILint(detect internal.ExecutableToolDetectFunc) *GolangCILint {
	return &GolangCILint{
		internal.MakeExecutableTool(
			IdGolangCILint,
			"Golangci-lint",
			"Smart, fast linters runner.",
			internal.Tags{internal.ToolTagGo, internal.ToolTagLint},
			detect,
		),
	}
}

// DetectGolangCILint create golangci-lint tool can be used in the project.
func DetectGolangCILint(
	ctx context.Context,
	options internal.DetectOptions,
) *GolangCILint {
	return NewGolangCILint(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(
			ctx,
			options.GetToolPath(IdGolangCILint, "golangci-lint"),
		)
	})
}

func (g *GolangCILint) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"run"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	if options.Output.Reporter != nil {
		return g.lintFilesReporter(ctx, options, options.Output.Reporter, args)
	}

	return g.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (g *GolangCILint) FormatFiles(
	ctx context.Context,
	options internal.ProjectFormatterOptions,
) error {
	args := append([]string{"fmt"}, options.ExtraArgs...)

	if options.CheckOnly {
		args = append(args, "--diff")
	}

	if options.Filenames != nil {
		args = append(args, *options.Filenames...)
	}

	return g.RunForProject(ctx, options.ProjectInfo, args)
}

func (g *GolangCILint) lintFilesReporter(
	ctx context.Context,
	options internal.ProjectLinterOptions,
	reporter internal.FileIssueReporter,
	args []string,
) error {
	// Process SARIF output
	args = append(args, "--output.sarif.path", "stdout")

	outputBuffer := bytes.Buffer{}

	err := g.RunForProjectWithWriter(ctx, options.ProjectInfo, &outputBuffer, args)

	if err != nil && !checkExitCode(err, []int{1}) {
		return err
	} else if outputBuffer.Len() == 0 {
		// Nothing returned
		return nil
	}

	return errors.Join(err, g.reportSarifOutput(&outputBuffer, reporter))
}

func (g *GolangCILint) reportSarifOutput(
	reader io.Reader,
	linterReporter internal.FileIssueReporter,
) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read JSON output: %w", err)
	}

	report, err := sarif.FromBytes(data)
	if err != nil {
		return fmt.Errorf("failed to parse SARIF output: %w", err)
	}

	return linterReporter.AddSarif(report)
}
