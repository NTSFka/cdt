package command

import (
	"context"
	"errors"
	"fmt"

	"cdt/internal"

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
	cmdContext := ctx.Value("context").(internal.Context)
	runner := cmdContext.Project.Workflow.Runner

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := cmdContext.Tools.Get(toolName)

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

	options := internal.ProjectRunnerOptions{
		ProjectInfo: cmdContext.Project.Info,
		Runtime:     cmdContext.Environment,
		ExtraArgs:   cmd.Args().Tail(),
	}

	if err := runner.RunTarget(ctx, options, target); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
