package command

import (
	"cdt/internal"
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewToolCommand() *cli.Command {
	return &cli.Command{
		Name:  "tool",
		Usage: "Work with supported tools",
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "List available tools",
				Action: toolCommandListAction,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "tag",
						Usage:   "List tools with specific tag",
						Aliases: []string{"t"},
					},
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Usage:   "List all supported tools",
					},
				},
			},
			{
				Name:   "run",
				Usage:  "Run a tool",
				Action: toolCommandRunAction,
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "toolId",
					},
				},
			},
		},
	}
}

func toolCommandListAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	var tools internal.Tools

	if cmd.Bool("all") {
		tools = c.Tools
	} else {
		tools = c.Tools.OnlyAvailable()
	}

	if cmd.IsSet("tag") {
		tags := cmd.StringSlice("tag")
		tools = tools.FilterByTags(tags)
	}

	tools.PrintTable(cmd.Writer)

	return nil
}

func toolCommandRunAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	toolId := cmd.StringArg("toolId")

	if toolId == "" {
		return fmt.Errorf("tool ID is required")
	}

	if tool := c.Tools.Get(toolId); tool != nil {
		options := internal.RunOptions{
			Directory: c.Project.Info.Directory,
			Input:     cmd.Reader,
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
