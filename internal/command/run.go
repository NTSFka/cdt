package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

func NewRunCommand() *cli.Command {
	return &cli.Command{
		Name:   "run",
		Usage:  "Run an application in the project",
		Action: runCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific run tool",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "target",
			},
		},
	}
}

func runCommandAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	runner := c.Project.Workflow.Runner

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := c.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		runnerTool, ok := tool.(internal.ProjectRunner)

		if ok {
			runner = runnerTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support run of target", toolName)
		}
	}

	if runner == nil {
		return errors.New("project doesn't support run of target")
	}

	target := cmd.StringArg("target")

	if target == "" {
		return errors.New("target is required")
	}

	if err := runner.RunTarget(c.Project, target, cmd.Args().Slice()); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
