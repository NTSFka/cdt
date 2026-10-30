package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

const IdPHPMD = "phpmd"

type PHPMDVersion int

const (
	PHPMDVersion2 PHPMDVersion = 2
	PHPMDVersion3 PHPMDVersion = 3
)

type PHPMD struct {
	internal.ExecutableTool
}

// DetectPHPMD create a tool for phpmd.
func DetectPHPMD(
	ctx context.Context,
	options internal.DetectOptions,
) *PHPMD {
	if execPath, ok := options.ToolsPaths[IdPHPMD]; ok {
		return NewPHPMD(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, execPath)
		})
	}

	return NewPHPMD(internal.DetectExecutableChain(
		[]string{
			"phpmd.phar",
			filepath.Join(options.ProjectDirectory, "vendor/bin/phpmd"),
			"phpmd",
		},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewPHPMD creates a phpmd tool from a custom executable.
func NewPHPMD(detect internal.ExecutableToolDetectFunc) *PHPMD {
	return &PHPMD{
		ExecutableTool: internal.MakeExecutableTool(
			IdPHPMD,
			"PHPMD",
			"PHPMD is a spin-off project of PHP Depend and aims to be a PHP equivalent of the well known Java tool PMD.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *PHPMD) AnalyzeFiles(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
) error {
	// Version 2 and 3 have different argument layout.
	version, err := p.detectVersion(ctx)
	if err != nil {
		return err
	}

	var args []string

	switch version {
	case PHPMDVersion2:
		args, err = p.argsBuild2(options)
	case PHPMDVersion3:
		args = p.argsBuild3(options)
		err = nil
	}

	if err != nil {
		return err
	}

	if options.Output.Reporter != nil {
		return p.analyzeFilesReporter(ctx, options, options.Output.Reporter, args)
	}

	return p.RunForProjectWithWriter(ctx, options.ProjectInfo, options.Output.Output, args)
}

func (p *PHPMD) analyzeFilesReporter(
	ctx context.Context,
	options internal.ProjectAnalyzerOptions,
	reporter internal.FileIssueReporter,
	args []string,
) error {
	// NOTE: The output format is already in args (version 2 vs. 3)
	outputBuffer := bytes.Buffer{}

	err := p.RunForProjectWithWriter(ctx, options.ProjectInfo, &outputBuffer, args)

	if err != nil && !checkExitCode(err, []int{2}) {
		return err
	} else if outputBuffer.Len() == 0 {
		// Nothing returned
		return nil
	}

	return errors.Join(err, p.reportSarifOutput(&outputBuffer, reporter))
}

func (p *PHPMD) reportSarifOutput(
	reader io.Reader,
	analyzerReporter internal.FileIssueReporter,
) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read JSON output: %w", err)
	}

	report, err := sarif.FromBytes(data)
	if err != nil {
		return fmt.Errorf("failed to parse SARIF output: %w", err)
	}

	return analyzerReporter.AddSarif(report)
}

func (p *PHPMD) detectVersion(ctx context.Context) (PHPMDVersion, error) {
	var output = bytes.Buffer{}

	if err := p.Run(ctx, internal.RunOptions{Output: &output}, []string{"--version"}); err != nil {
		return 0, err
	}

	if strings.Contains(output.String(), "PHPMD 2") {
		return PHPMDVersion2, nil
	}

	if strings.Contains(output.String(), "PHPMD 3") {
		return PHPMDVersion3, nil
	}

	return 0, fmt.Errorf("unsupported PHPMD version: %s", output.String())
}

func (p *PHPMD) argsBuild2(options internal.ProjectAnalyzerOptions) ([]string, error) {
	var args []string

	args = append(args, p.argsBuild2Fileset(options))

	if options.Output.Reporter != nil {
		args = append(args, "sarif")
	} else {
		args = append(args, "ansi")
	}

	if a, err := p.argsBuild2Ruleset(options); err == nil {
		args = append(args, a...)
	} else {
		return nil, err
	}

	args = append(args, options.ExtraArgs...)

	return args, nil
}

func (p *PHPMD) argsBuild3(options internal.ProjectAnalyzerOptions) []string {
	args := []string{"analyze"}

	if options.Output.Reporter != nil {
		args = append(args, "--format=sarif")
	}

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	}

	return append(args, options.ExtraArgs...)
}

func (p *PHPMD) argsBuild2Fileset(options internal.ProjectAnalyzerOptions) string {
	if options.Filenames != nil && len(*options.Filenames) > 0 {
		return strings.Join(*options.Filenames, ",")
	}

	// Try most common names for directories
	var paths []string

	for _, name := range []string{"src", "app", "bin", "tests"} {
		if internal.PathExists(path.Join(options.Directory, name)) {
			paths = append(paths, name)
		}
	}

	if len(paths) > 0 {
		return strings.Join(paths, ",")
	}

	// Everything in the project directory
	return options.Directory
}

func (p *PHPMD) argsBuild2Ruleset(options internal.ProjectAnalyzerOptions) ([]string, error) {
	var args []string

	switch {
	case internal.PathExists(path.Join(options.Directory, "phpmd.xml")):
		args = append(args, "phpmd.xml")
	case internal.PathExists(path.Join(options.Directory, "phpmd.ruleset.xml")):
		args = append(args, "phpmd.ruleset.xml")
	case len(options.ExtraArgs) == 0:
		return nil, fmt.Errorf(
			"ruleset required, create a phpmd.xml or pass additional argument for ruleset",
		)
	}

	return args, nil
}
