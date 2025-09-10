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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific configurator tool",
			},
		},
	}
}

func configureCommandAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	configurator := c.Workflow.Configurator

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := c.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		configuratorTool, ok := tool.(internal.ProjectConfigurator)

		if ok {
			configurator = configuratorTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support configuration", toolName)
		}
	}

	if configurator == nil {
		return errors.New("project doesn't support configuration")
	}

	if err := configurator.Configure(c.ProjectInfo, cmd.Args().Tail()); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
