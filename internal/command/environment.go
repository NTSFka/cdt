package command

import (
	"cdt/internal"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
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

func environmentCommandActionList(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	c.EnvironmentProviders.PrintList(os.Stdout)

	return nil
}

func environmentCommandActionStatus(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	env := c.Environment

	if env.IsRunning() {
		fmt.Printf("%v: running\n", env.Id())
	} else {
		fmt.Printf("%v: stopped\n", env.Id())
	}

	return nil
}

func environmentCommandActionStart(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	env := c.Environment

	if err := env.Start(); err != nil {
		return fmt.Errorf("environment start failed: %w", err)
	}

	return nil
}

func environmentCommandActionStop(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	env := c.Environment

	if err := env.Stop(); err != nil {
		return fmt.Errorf("environment stop failed: %w", err)
	}

	return nil
}
