package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"encoding/json"
	"errors"
	"io"
	"path/filepath"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

const IdPHPStan = "phpstan"

type PHPStan struct {
	internal.ExecutableTool
}

// DetectPHPStan create a tool for phpstan.
func DetectPHPStan(
	ctx context.Context,
	options internal.DetectOptions,
) *PHPStan {
	if path, ok := options.ToolsPaths[IdPHPStan]; ok {
		return NewPHPStan(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewPHPStan(internal.DetectExecutableChain(
		[]string{filepath.Join(options.ProjectDirectory, "vendor/bin/phpstan"), "phpstan"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPStan creates a phpstan tool from a custom executable.
func NewPHPStan(detect internal.ExecutableToolDetectFunc) *PHPStan {
	return &PHPStan{
		ExecutableTool: internal.MakeExecutableTool(
			IdPHPStan,
			"PHPStan",
			"Analyses source code.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagAnalyze},
			detect,
		),
	}
}

func (p *PHPStan) AnalyzeFiles(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
) error {
	args := append([]string{"analyse"}, options.ExtraArgs...)
	args = appendFiles(args, options.Filenames, nil)

	if options.Output.Reporter != nil {
		return p.analyzeFilesReporter(ctx, options, options.Output.Reporter, args)
	}

	return p.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (p *PHPStan) analyzeFilesReporter(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
	reporter internal.FileIssueReporter,
	args []string,
) error {
	// Process JSON output
	args = append(args, "--error-format=json", "--quiet")

	outputBuffer := bytes.Buffer{}

	err := p.RunForProjectWithWriter(ctx, options.ProjectInfo, &outputBuffer, args)

	if err != nil && !checkExitCode(err, []int{1}) {
		return err
	} else if outputBuffer.Len() == 0 {
		// Nothing returned
		return nil
	}

	return errors.Join(err, p.reportJsonOutput(&outputBuffer, reporter))
}

type phpStanJsonReport struct {
	// Overview
	Totals struct {
		Errors     int `json:"errors"`
		FileErrors int `json:"file_errors"`
	} `json:"totals"`

	// Files issues
	Files map[string]struct {
		Errors   int `json:"errors"`
		Messages []struct {
			Message    string `json:"message"`
			Line       int    `json:"line"`
			Ignorable  bool   `json:"ignorable"`
			Identifier string `json:"identifier"`
		} `json:"messages"`
	} `json:"files"`

	// Other errors
	Errors []string `json:"errors"`
}

func (p *PHPStan) reportJsonOutput(
	reader io.Reader,
	analyzerReporter internal.FileIssueReporter,
) error {
	var report phpStanJsonReport
	if err := json.NewDecoder(reader).Decode(&report); err != nil {
		return err
	}

	sarifReport := sarif.NewReport()

	run := sarif.NewRunWithInformationURI(
		p.Name(),
		"https://phpstan.org",
	)

	for path, info := range report.Files {
		for _, message := range info.Messages {
			run.CreateResultForRule(message.Identifier).
				WithMessage(sarif.NewTextMessage(message.Message)).
				AddLocation(sarif.NewLocationWithPhysicalLocation(
					sarif.NewPhysicalLocation().
						WithArtifactLocation(sarif.NewSimpleArtifactLocation(path)).
						WithRegion(sarif.NewRegion().WithStartLine(message.Line)),
				))
		}
	}

	sarifReport.AddRun(run)

	return analyzerReporter.AddSarif(sarifReport)
}
