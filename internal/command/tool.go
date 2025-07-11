package command

import (
	"cdt/internal"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
)

var ToolCommand = cli.Command{
	Name:  "tool",
	Usage: "work with supported tools",
	Commands: []*cli.Command{
		{
			Name:   "list",
			Usage:  "list available tools",
			Action: toolCommandListAction,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "all",
					Aliases: []string{"a"},
					Usage:   "list all supported tools",
				},
			},
		},
		{
			Name:   "run",
			Usage:  "run a tool",
			Action: toolCommandRunAction,
			Arguments: []cli.Argument{
				&cli.StringArg{
					Name: "toolId",
				},
			},
		},
	},
}

func toolCommandListAction(ctx context.Context, command *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	if command.Bool("all") {
		internal.PrintToolList(os.Stdout, c.Tools)
	} else {
		internal.PrintToolList(os.Stdout, c.Tools.Active())
	}

	return nil
}

func toolCommandRunAction(ctx context.Context, command *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	toolId := command.StringArg("toolId")

	if tool := c.Tools.Get(toolId); tool != nil {
		err := tool.Run(c.Project, command.Args().Slice())

		if err != nil {
			return fmt.Errorf("tool '%s' failed: %w", toolId, err)
		}

		return nil
	}

	return fmt.Errorf("tool '%s' not found", toolId)
}
