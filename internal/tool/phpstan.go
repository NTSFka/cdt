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

// nolint: cyclop
func (p *PHPStan) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"analyse"}, options.ExtraArgs...)

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	}

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.LintReportFormatDefault:
		fallthrough
	case internal.LintReportFormatRaw:
		break
	case internal.LintReportFormatJson:
		args = append(args, "--error-format=json")
	case internal.LintReportFormatJUnit:
		args = append(args, "--error-format=junit")
	case internal.LintReportFormatGitHub:
		args = append(args, "--error-format=github")
	case internal.LintReportFormatGitLab:
		args = append(args, "--error-format=gitlab")
	case internal.LintReportFormatTeamCity:
		args = append(args, "--error-format=teamcity")
	default:
		return fmt.Errorf("unsupported report format: %s", options.Output.Format)
	}

	if options.Output.Filename != nil {
		return p.RunForProjectWithOutput(ctx, options.ProjectInfo, *options.Output.Filename, args)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
