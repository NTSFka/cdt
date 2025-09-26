package command

import (
	"context"
	"errors"
	"fmt"

	"cdt/internal"

	"github.com/urfave/cli/v3"
)

func NewTestCommand() *cli.Command {
	return &cli.Command{
		Name:   "test",
		Usage:  "TestPattern the project",
		Action: testCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific test tool",
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

	options := internal.ProjectTesterOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	var err error

	if pattern := cmd.StringArgs("pattern"); len(pattern) != 0 {
		err = tester.TestPattern(ctx, options, pattern[0])
	} else {
		err = tester.TestAll(ctx, options)
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
