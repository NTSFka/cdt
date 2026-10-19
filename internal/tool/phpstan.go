package tool

import (
	"cdt/internal"
	"context"
	"fmt"
	"path/filepath"
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

	if a, err := p.argsBuildLintOutputFormat(options.Output.Format); err == nil {
		args = append(args, a...)
	} else {
		return err
	}

	if options.Output.Filename != nil {
		return p.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
}

func (p *PHPStan) argsBuildLintOutputFormat(format internal.AnalyzeReportFormat) ([]string, error) {
	var args []string

	// nolint: exhaustive
	switch format {
	case internal.AnalyzeReportFormatDefault:
		fallthrough
	case internal.AnalyzeReportFormatRaw:
		break
	case internal.AnalyzeReportFormatJson:
		args = []string{"--error-format=json"}
	case internal.AnalyzeReportFormatJUnit:
		args = []string{"--error-format=junit"}
	case internal.AnalyzeReportFormatGitHub:
		args = []string{"--error-format=github"}
	case internal.AnalyzeReportFormatGitLab:
		args = []string{"--error-format=gitlab"}
	case internal.AnalyzeReportFormatTeamCity:
		args = []string{"--error-format=teamcity"}
	default:
		return nil, fmt.Errorf("unsupported report format: %s", format)
	}

	return args, nil
}
