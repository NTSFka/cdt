package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

var RunCommand = cli.Command{
	Name:      "run",
	Usage:     "Run an application in the project",
	UsageText: "cdt run TARGET",
	Action:    runCommandAction,
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name: "target",
		},
	},
}

func runCommandAction(ctx context.Context, command *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	runner := c.Project.Runner()

	if runner == nil {
		return errors.New("project doesn't support run of target")
	}

	target := command.StringArg("target")

	if target == "" {
		return errors.New("target is required")
	}

	if err := runner.Run(target, command.Args().Slice()); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
