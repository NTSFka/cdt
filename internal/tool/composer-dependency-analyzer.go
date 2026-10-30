package tool

import (
	"bytes"
	"cdt/internal"
	"cdt/internal/utils/junit"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

const IdComposerDependencyAnalyzer = "composer-dependency-analyzer"

type ComposerDependencyAnalyzer struct {
	internal.ExecutableTool
}

// DetectComposerDependencyAnalyzer create a tool for composer-dependency-analyzer.
func DetectComposerDependencyAnalyzer(
	ctx context.Context,
	options internal.DetectOptions,
) *ComposerDependencyAnalyzer {
	if path, ok := options.ToolsPaths[IdComposerDependencyAnalyzer]; ok {
		return NewComposerDependencyAnalyzer(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewComposerDependencyAnalyzer(internal.DetectExecutableChain(
		[]string{
			filepath.Join(options.ProjectDirectory, "vendor/bin/composer-dependency-analyzer"),
			"composer-dependency-analyzer",
		},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewComposerDependencyAnalyzer creates a composer-dependency-analyzer tool from a custom executable.
func NewComposerDependencyAnalyzer(
	detect internal.ExecutableToolDetectFunc,
) *ComposerDependencyAnalyzer {
	return &ComposerDependencyAnalyzer{
		ExecutableTool: internal.MakeExecutableTool(
			IdComposerDependencyAnalyzer,
			"ComposerDependencyAnalyzer",
			"Detects unused, shadow and misplaced composer dependencies.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagAnalyze},
			detect,
		),
	}
}

func (c *ComposerDependencyAnalyzer) AnalyzeFiles(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
) error {
	args := options.ExtraArgs
	args = appendFiles(args, options.Filenames, nil)

	if options.Output.Reporter != nil {
		return c.analyzeFilesReporter(ctx, options, options.Output.Reporter, args)
	}

	return c.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (c *ComposerDependencyAnalyzer) analyzeFilesReporter(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
	reporter internal.FileIssueReporter,
	args []string,
) error {
	// Process JSON output
	args = append(args, "--format", "junit")

	outputBuffer := bytes.Buffer{}

	err := c.RunForProjectWithWriter(ctx, options.ProjectInfo, &outputBuffer, args)

	if err != nil && !checkExitCode(err, []int{1}) {
		return err
	} else if outputBuffer.Len() == 0 {
		// Nothing returned
		return nil
	}

	return errors.Join(err, c.reportJUnitOutput(&outputBuffer, reporter))
}

func (c *ComposerDependencyAnalyzer) reportJUnitOutput(
	reader io.Reader,
	analyzerReporter internal.FileIssueReporter,
) error {
	var report junit.Testsuites
	if err := xml.NewDecoder(reader).Decode(&report); err != nil {
		// If there is no issue, it returns XML without root element
		//  Sadly, if the XML is not valid, it returns the same error
		if errors.Is(err, io.EOF) {
			return nil
		}

		return err
	}

	sarifReport := sarif.NewReport()

	run := sarif.NewRunWithInformationURI(
		c.Name(),
		"https://github.com/shipmonk-rnd/composer-dependency-analyser",
	)

	for _, suite := range report.Suites {
		for _, testcase := range suite.Testcases {
			failure := testcase.Failure

			location := strings.SplitN(failure.Data, ":", 2)

			if len(location) != 2 {
				return fmt.Errorf("invalid data: %s", failure.Data)
			}

			line, _ := strconv.Atoi(location[1])

			run.CreateResultForRule("").
				WithMessage(sarif.NewTextMessage(failure.Message)).
				AddLocation(sarif.NewLocationWithPhysicalLocation(
					sarif.NewPhysicalLocation().
						WithArtifactLocation(sarif.NewSimpleArtifactLocation(location[0])).
						WithRegion(sarif.NewRegion().WithStartLine(line)),
				))
		}
	}

	sarifReport.AddRun(run)

	return analyzerReporter.AddSarif(sarifReport)
}
