package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewFormatCommand() *cli.Command {
	return &cli.Command{
		Name:   "format",
		Usage:  "Format files in the project",
		Action: formatCommandAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "check",
				Value: false,
				Usage: "Check if the project or given files are formatted",
			},
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific formatter tool",
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

func formatCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	formatter := cmdContext.Project.Workflow.Formatter

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := cmdContext.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		formatterTool, ok := tool.(internal.ProjectFormatter)

		if ok {
			formatter = formatterTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support formatting", toolName)
		}
	}

	if formatter == nil {
		return errors.New("project doesn't support source formatting")
	}

	options := internal.ProjectFormatterOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	var err error

	if cmd.Bool("check") { // nolint: nestif
		if files := cmd.StringArgs("files"); len(files) > 0 {
			err = formatter.FormatCheckFiles(ctx, options, files)
		} else {
			err = formatter.FormatCheckAll(ctx, options)
		}
	} else {
		if files := cmd.StringArgs("files"); len(files) > 0 {
			err = formatter.FormatFiles(ctx, options, files)
		} else {
			err = formatter.FormatAll(ctx, options)
		}
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
