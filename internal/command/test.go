package command

import (
	"cdt/internal"
	"cdt/internal/reporter"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/urfave/cli/v3"
)

func NewTestCommand() *cli.Command {
	return &cli.Command{
		Name:   "test",
		Usage:  "Test the project",
		Action: testCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific test tool",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage: "The test report output format in form: <format>[:<filename>]. When <filename> is not " +
					"specified the result is written to stdout. Available formats: raw, json, ctrf",
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

func testCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	tester := cmdContext.Project.Workflow.Tester

	if cmd.IsSet("tool") {
		var err error

		tester, err = getTool[internal.ProjectTester](
			cmd.String("tool"),
			cmdContext.Tools,
			"testing",
		)
		if err != nil {
			return err
		}
	}

	if tester == nil {
		return errors.New("project doesn't support testing")
	}

	outputOptions := ParseOptionOutput(cmd.String("output"))

	output, err := buildOutput(outputOptions.Filename)
	if err != nil {
		return err
	}
	defer func() { _ = output.Close() }()

	rep, err := buildTestsReporter(outputOptions.Format, output)
	if err != nil {
		return err
	}

	options := internal.ProjectTesterOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
		Output: internal.TestsOutputOptions{
			Output:   output,
			Reporter: rep,
		},
	}

	if pattern := cmd.StringArgs("pattern"); len(pattern) != 0 {
		options.Pattern = &pattern[0]
	}

	if err := tester.RunTests(ctx, options); err != nil {
		return fmt.Errorf("%w", err)
	}

	if rep != nil {
		return rep.Finish()
	}

	return nil
}

func buildTestsReporter(format string, output io.Writer) (internal.TestsReporter, error) {
	switch format {
	case "":
		fallthrough
	case "raw":
		return nil, nil // nolint: nilnil
	case "simple":
		return reporter.NewSimpleCtrfReporter(output), nil
	case "sarif":
		return reporter.NewFileCtrfReporter(output), nil
	}

	return nil, fmt.Errorf("unsupported report format: %s", format)
}
