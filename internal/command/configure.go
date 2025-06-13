package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

var ConfigureCommand = cli.Command{
	Name:      "configure",
	Usage:     "Configure the project",
	UsageText: "cdt [OPTIONS] configure",
	Action:    configureCommandAction,
}

func configureCommandAction(ctx context.Context, command *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	configurator := c.Workflow.Configurator

	if configurator == nil {
		return errors.New("project doesn't support configuration")
	}

	if err := configurator.Configure(c.Project, command.Args().Tail()); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
