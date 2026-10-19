package tool

import (
	"bytes"
	"cdt/internal"
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
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
		args, err = p.argsBuild3(options)
	}

	if err != nil {
		return err
	}

	if options.Output.Filename != nil {
		return p.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
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

	if arg, err := p.argsBuildAnalyzeOutputFormat(options.Output.Format, new("ansi")); err == nil {
		args = append(args, *arg)
	} else {
		return nil, err
	}

	if a, err := p.argsBuild2Ruleset(options); err == nil {
		args = append(args, a...)
	} else {
		return nil, err
	}

	args = append(args, options.ExtraArgs...)

	return args, nil
}

func (p *PHPMD) argsBuild3(options internal.ProjectAnalyzerOptions) ([]string, error) {
	args := []string{"analyze"}

	if arg, err := p.argsBuildAnalyzeOutputFormat(options.Output.Format, nil); err == nil {
		if arg != nil {
			args = append(args, "--format="+*arg)
		}
	} else {
		return nil, err
	}

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	}

	args = append(args, options.ExtraArgs...)

	return args, nil
}

func (p *PHPMD) argsBuildAnalyzeOutputFormat(
	format internal.AnalyzeReportFormat,
	def *string,
) (*string, error) {
	// supported: ansi, baseline, checkstyle, github, gitlab, html, json, sarif, text, xml.

	// nolint: exhaustive
	switch format {
	case internal.AnalyzeReportFormatDefault:
		fallthrough
	case internal.AnalyzeReportFormatRaw:
		return def, nil
	case internal.AnalyzeReportFormatJson:
		return new("json"), nil
	case internal.AnalyzeReportFormatGitHub:
		return new("github"), nil
	case internal.AnalyzeReportFormatGitLab:
		return new("gitlab"), nil
	case internal.AnalyzeReportFormatXml:
		return new("xml"), nil
	case internal.AnalyzeReportFormatHtml:
		return new("html"), nil
	case internal.AnalyzeReportFormatSarif:
		return new("sarif"), nil
	default:
		return nil, fmt.Errorf("unsupported report format: %s", format)
	}
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
