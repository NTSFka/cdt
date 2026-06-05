package command

import (
	"context"
	"errors"
	"fmt"

	"cdt/internal"

	"github.com/urfave/cli/v3"
)

func NewCoverageCommand() *cli.Command {
	return &cli.Command{
		Name:    "coverage",
		Usage:   "run tests with coverage for the project",
		Aliases: []string{"cov"},
		Action:  coverageCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific test tool",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "The coverage report output format in form: <format>[:<filename/directory>]",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "pattern",
				Min:  0,
				Max:  1,
			},
		},
	}
}

func coverageCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	collector := cmdContext.Project.Workflow.CoverageCollector

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := cmdContext.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		collectorTool, ok := tool.(internal.ProjectCoverageCollector)

		if ok {
			collector = collectorTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support coverage collection", toolName)
		}
	}

	if collector == nil {
		return errors.New("project doesn't support coverage collection")
	}

	outputOptions := ParseOptionOutput[internal.CoverageReportFormat](
		cmd.String("output"),
		internal.CoverageReportFormatRaw,
	)

	options := internal.ProjectCoverageCollectorOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
		Output:      outputOptions,
	}

	if pattern := cmd.StringArgs("pattern"); len(pattern) != 0 {
		options.Pattern = &pattern[0]
	}

	err := collector.CollectCoverage(ctx, options)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
