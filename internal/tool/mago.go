package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

const IdMago = "mago"

type Mago struct {
	internal.ExecutableTool
}

// DetectMago create a tool for mago.
func DetectMago(
	ctx context.Context,
	options internal.DetectOptions,
) *Mago {
	if path, ok := options.ToolsPaths[IdMago]; ok {
		return NewMago(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewMago(internal.DetectExecutableChain(
		[]string{filepath.Join(options.ProjectDirectory, "vendor/bin/mago"), "mago"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewMago creates a mago tool from a custom executable.
func NewMago(detect internal.ExecutableToolDetectFunc) *Mago {
	return &Mago{
		ExecutableTool: internal.MakeExecutableTool(
			IdMago,
			"Mago",
			"Mago is a static analyzer, linter, and formatter for PHP, written in Rust.",
			internal.Tags{
				internal.ToolTagPhp,
				internal.ToolTagLint,
				internal.ToolTagAnalyze,
				internal.ToolTagFormat,
			},
			detect,
		),
	}
}

func (m *Mago) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"lint"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	if options.Output.Reporter != nil {
		return m.runForSarifReporter(ctx, options.ProjectInfo, options.Output.Reporter, args)
	}

	return m.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (m *Mago) AnalyzeFiles(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
) error {
	args := append([]string{"analyze"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	if options.Output.Reporter != nil {
		return m.runForSarifReporter(ctx, options.ProjectInfo, options.Output.Reporter, args)
	}

	return m.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (m *Mago) FormatFiles(ctx context.Context, options internal.ProjectFormatterOptions) error {
	args := append([]string{"format"}, options.ExtraArgs...)

	if options.CheckOnly {
		args = append(args, "--dry-run")
	}

	args = appendFiles(args, options.Filenames, nil)

	return m.RunForProject(ctx, options.ProjectInfo, args)
}

func (m *Mago) runForSarifReporter(
	ctx context.Context,
	projectInfo internal.ProjectInfo,
	reporter internal.FileIssueReporter,
	args []string,
) error {
	// Process SARIF output
	args = append(args, "--reporting-format=sarif")

	outputBuffer := bytes.Buffer{}

	err := m.RunForProjectWithWriter(ctx, projectInfo, &outputBuffer, args)

	if err != nil && !checkExitCode(err, []int{1}) {
		return err
	} else if outputBuffer.Len() == 0 {
		// Nothing returned
		return nil
	}

	return errors.Join(err, m.reportSarifOutput(&outputBuffer, reporter))
}

func (m *Mago) reportSarifOutput(
	reader io.Reader,
	sarifReporter internal.FileIssueReporter,
) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read JSON output: %w", err)
	}

	report, err := sarif.FromBytes(data)
	if err != nil {
		return fmt.Errorf("failed to parse SARIF output: %w", err)
	}

	return sarifReporter.AddSarif(report)
}
