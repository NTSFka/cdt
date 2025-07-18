package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

func NewConfigureCommand() *cli.Command {
	return &cli.Command{
		Name:   "configure",
		Usage:  "Configure the project",
		Action: configureCommandAction,
	}
}

func configureCommandAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	configurator := c.Project.Workflow.Configurator

	if configurator == nil {
		return errors.New("project doesn't support configuration")
	}

	if err := configurator.Configure(c.Project, cmd.Args().Tail()); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
