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
	options DetectOptions,
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
			internal.Tags{internal.ToolTagPhp, internal.ToolTagLint},
			detect,
		),
	}
}

func (p *PHPStan) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
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

func (p *PHPStan) argsBuildLintOutputFormat(format internal.LintReportFormat) ([]string, error) {
	var args []string

	switch format {
	case internal.LintReportFormatDefault:
		fallthrough
	case internal.LintReportFormatRaw:
		break
	case internal.LintReportFormatJson:
		args = []string{"--error-format=json"}
	case internal.LintReportFormatJUnit:
		args = []string{"--error-format=junit"}
	case internal.LintReportFormatGitHub:
		args = []string{"--error-format=github"}
	case internal.LintReportFormatGitLab:
		args = []string{"--error-format=gitlab"}
	case internal.LintReportFormatTeamCity:
		args = []string{"--error-format=teamcity"}
	default:
		return nil, fmt.Errorf("unsupported report format: %s", format)
	}

	return args, nil
}
