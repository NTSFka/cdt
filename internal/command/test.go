package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

var TestCommand = cli.Command{
	Name:      "test",
	Usage:     "Test the project",
	Action:    testCommandAction,
	UsageText: "cdt [OPTIONS] test [PATTERNS...]",
	Arguments: []cli.Argument{
		&cli.StringArgs{
			Name: "pattern",
			Min:  0,
			Max:  1,
		},
	},
}

func testCommandAction(ctx context.Context, command *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	tester := c.Workflow.Tester

	if tester == nil {
		return errors.New("project doesn't support testing")
	}

	var err error

	if pattern := command.StringArgs("pattern"); len(pattern) != 0 {
		err = tester.Test(c.Project, pattern[0])
	} else {
		err = tester.TestAll(c.Project)
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
