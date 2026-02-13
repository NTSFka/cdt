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

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	}

	// nolint: exhaustive
	switch options.Output.Format {
	case internal.LintReportFormatDefault:
		fallthrough
	case internal.LintReportFormatRaw:
		break
	default:
		return fmt.Errorf("unsupported report format: %s", options.Output.Format)
	}

	return p.RunForProject(ctx, options.ProjectInfo, args)
}
