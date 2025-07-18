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
	c := ctx.Value("context").(internal.Context)
	tester := c.Project.Workflow.Tester

	if tester == nil {
		return errors.New("project doesn't support testing")
	}

	var err error

	if pattern := cmd.StringArgs("pattern"); len(pattern) != 0 {
		err = tester.Test(c.Project, pattern[0], cmd.Args().Tail())
	} else {
		err = tester.TestAll(c.Project, cmd.Args().Tail())
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
