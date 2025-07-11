package command

import (
	"cdt/internal"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

var EnvironmentCommand = cli.Command{
	Name:    "environment",
	Aliases: []string{"env"},
	Usage:   "manage runtime environment",
	Commands: []*cli.Command{
		{
			Name:   "list",
			Action: environmentCommandActionList,
		},
		{
			Name:   "status",
			Action: environmentCommandActionStatus,
		},
		{
			Name:   "start",
			Action: environmentCommandActionStart,
		},
		{
			Name:   "stop",
			Action: environmentCommandActionStop,
		},
	},
}

func environmentCommandActionList(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	c.EnvironmentProviders.PrintTable(cmd.Writer)

	return nil
}

func environmentCommandActionStatus(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	env := c.Environment

	if env.IsRunning(ctx) {
		fmt.Fprintf(cmd.Writer, "%v: running\n", env.Id())
	} else {
		fmt.Fprintf(cmd.Writer, "%v: stopped\n", env.Id())
	}

	return nil
}

func environmentCommandActionStart(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	env := c.Environment

	if err := env.Start(ctx); err != nil {
		return fmt.Errorf("environment start failed: %w", err)
	}

	return nil
}

func environmentCommandActionStop(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	env := c.Environment

	if err := env.Stop(ctx); err != nil {
		return fmt.Errorf("environment stop failed: %w", err)
	}

	return nil
}
