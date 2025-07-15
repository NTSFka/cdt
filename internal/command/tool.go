package command

import (
	"cdt/internal"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
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

func toolCommandListAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	if cmd.Bool("all") {
		c.Tools.PrintTable(cmd.Writer)
	} else {
		tools := c.Tools.OnlyAvailable()
		tools.PrintTable(cmd.Writer)
	}

	return nil
}

func toolCommandRunAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	toolId := cmd.StringArg("toolId")

	if tool := c.Tools.Get(toolId); tool != nil {
		options := internal.RunOptions{
			Directory: c.Project.RootDirectory(),
			Output:    cmd.Writer,
			Error:     cmd.ErrWriter,
		}

		err := tool.Run(ctx, options, cmd.Args().Slice())

		if err != nil {
			return fmt.Errorf("tool '%s' failed: %w", toolId, err)
		}

		return nil
	}

	return fmt.Errorf("tool '%s' not found", toolId)
}
