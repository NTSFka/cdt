package command

import (
	"cdt/internal"
	"context"
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
