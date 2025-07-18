package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

func NewRunCommand() *cli.Command {
	return &cli.Command{
		Name:   "run",
		Usage:  "Run an application in the project",
		Action: runCommandAction,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "target",
			},
		},
	}
}

func runCommandAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	runner := c.Project.Workflow.Runner

	if runner == nil {
		return errors.New("project doesn't support run of target")
	}

	target := cmd.StringArg("target")

	if target == "" {
		return errors.New("target is required")
	}

	if err := runner.RunTarget(c.Project, target, cmd.Args().Slice()); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
