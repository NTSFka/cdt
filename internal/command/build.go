package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

var BuildCommand = cli.Command{
	Name:   "build",
	Usage:  "build the whole project or target(s) in the project",
	Action: buildCommandAction,
	Arguments: []cli.Argument{
		&cli.StringArgs{
			Name: "targets",
			Min:  0,
			Max:  -1,
		},
	},
}

func buildCommandAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	builder := c.Project.Workflow.Builder

	if builder == nil {
		return errors.New("project doesn't support building")
	}

	var err error

	if targets := cmd.StringArgs("targets"); len(targets) > 0 {
		err = builder.BuildTargets(c.Project, targets, cmd.Args().Tail())
	} else {
		err = builder.BuildAll(c.Project, cmd.Args().Tail())
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
