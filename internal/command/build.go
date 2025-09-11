package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewBuildCommand() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "build the whole project or target(s) in the project",
		Action: buildCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific build tool",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "targets",
				Min:  0,
				Max:  -1,
			},
		},
	}
}

func buildCommandAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	builder := c.Project.Workflow.Builder

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := c.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		builderTool, ok := tool.(internal.ProjectBuilder)

		if ok {
			builder = builderTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support building", toolName)
		}
	}

	if builder == nil {
		return errors.New("project doesn't support building")
	}

	options := internal.ProjectBuilderOptions{
		ProjectInfo: c.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	var err error

	if targets := cmd.StringArgs("targets"); len(targets) > 0 {
		err = builder.BuildTargets(ctx, options, targets)
	} else {
		err = builder.BuildAll(ctx, options)
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
