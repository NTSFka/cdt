package tool

import (
	"context"
	"fmt"

	"cdt/internal"
)

const IdPHP = "php"

type PHP struct {
	internal.ExecutableTool
}

// DetectPHP create a tool for php.
func DetectPHP(
	ctx context.Context,
	options DetectOptions,
) *PHP {
	return NewPHP(func() (*internal.Executable, error) {
		return options.Environment.FindExecutable(ctx, options.GetToolPath(IdPHP, "php"))
	})
}

// NewPHP creates a php tool from a custom executable.
func NewPHP(detect internal.ExecutableToolDetectFunc) *PHP {
	return &PHP{
		ExecutableTool: internal.MakeExecutableTool(
			IdPHP,
			"PHP",
			"Popular general-purpose scripting language that is especially suited to web development.",
			internal.Tags{internal.ToolTagPhp, internal.ToolTagRun},
			detect,
		),
	}
}

func (p *PHP) RunTarget(
	ctx context.Context,
	options internal.ProjectRunnerOptions,
	target string,
) error {
	return p.RunForProject(
		ctx,
		options.ProjectInfo,
		append([]string{"-f", target}, options.ExtraArgs...),
	)
}

func (p *PHP) LintFiles(
	ctx context.Context,
	options internal.ProjectLinterOptions,
) error {
	args := append([]string{"-l"}, options.ExtraArgs...)

	if options.Filenames != nil && len(*options.Filenames) > 0 {
		args = append(args, *options.Filenames...)
	} else {
		return fmt.Errorf("no files to lint")
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
