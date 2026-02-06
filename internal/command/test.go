package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"

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
		toolName := cmd.String("tool")
		tool := cmdContext.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		testerTool, ok := tool.(internal.ProjectTester)

		if ok {
			tester = testerTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support testing", toolName)
		}
	}

	if tester == nil {
		return errors.New("project doesn't support testing")
	}

	outputOptions := ParseOptionOutput[internal.TestsReportFormat](
		cmd.String("output"),
		internal.TestsReportFormatRaw,
	)

	options := internal.ProjectTesterOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
		Output:      outputOptions,
	}

	if pattern := cmd.StringArgs("pattern"); len(pattern) != 0 {
		options.Pattern = &pattern[0]
	}

	err := tester.RunTests(ctx, options)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
