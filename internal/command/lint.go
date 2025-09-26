package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewLintCommand() *cli.Command {
	return &cli.Command{
		Name:   "lint",
		Usage:  "Lint the project",
		Action: lintCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific linter tool",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "files",
				Min:  0,
				Max:  -1,
			},
		},
	}
}

func lintCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	linter := cmdContext.Project.Workflow.Linter

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := cmdContext.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		linterTool, ok := tool.(internal.ProjectLinter)

		if ok {
			linter = linterTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support linting", toolName)
		}
	}

	if linter == nil {
		return errors.New("project doesn't support linting")
	}

	options := internal.ProjectLinterOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	var err error
	if files := cmd.StringArgs("files"); len(files) > 0 {
		err = linter.LintFiles(ctx, options, files)
	} else {
		err = linter.LintAll(ctx, options)
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
