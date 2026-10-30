package command

import (
	"cdt/internal/reporter"
	"context"
	"errors"
	"fmt"
	"io"

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
		var err error

		collector, err = getTool[internal.ProjectCoverageCollector](
			cmd.String("tool"),
			cmdContext.Tools,
			"coverage collection",
		)
		if err != nil {
			return err
		}
	}

	if collector == nil {
		return errors.New("project doesn't support coverage collection")
	}

	outputOptions := ParseOptionOutput(cmd.String("output"))

	output, err := buildOutput(outputOptions.Filename)
	if err != nil {
		return err
	}
	defer func() { _ = output.Close() }()

	rep, err := buildCoverageReporter(outputOptions.Format, output)
	if err != nil {
		return err
	}

	options := internal.ProjectCoverageCollectorOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
		Output: internal.CoverageOutputOptions{
			Output:   output,
			Reporter: rep,
		},
	}

	if pattern := cmd.StringArgs("pattern"); len(pattern) != 0 {
		options.Pattern = &pattern[0]
	}

	if err := collector.CollectCoverage(ctx, options); err != nil {
		return fmt.Errorf("%w", err)
	}

	if rep != nil {
		return rep.Finish()
	}

	return nil
}

func buildCoverageReporter(format string, output io.Writer) (internal.CoverageReporter, error) {
	switch format {
	case "":
		fallthrough
	case "raw":
		return nil, nil // nolint: nilnil
	case "simple":
		return reporter.NewSimpleCoberturaReporter(output), nil
	}

	return nil, fmt.Errorf("unsupported report format: %s", format)
}
