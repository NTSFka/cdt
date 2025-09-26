package command

import (
	"context"
	"errors"
	"fmt"

	"cdt/internal"

	"github.com/urfave/cli/v3"
)

func NewExecCommand() *cli.Command {
	return &cli.Command{
		Name:   "exec",
		Usage:  "Execute a command in the environment",
		Action: execCommandAction,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "command",
			},
			&cli.StringArgs{
				Name: "args",
				Min:  0,
				Max:  -1,
			},
		},
	}
}

func execCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)

	command := cmd.StringArg("command")

	if command == "" {
		return errors.New("COMMAND is required")
	}

	options := internal.RunOptions{
		Directory: cmdContext.Project.Info.Directory,
		Input:     cmd.Reader,
		Output:    cmd.Writer,
		Error:     cmd.ErrWriter,
	}

	err := cmdContext.Environment.RunExecutable(ctx, options, command, cmd.StringArgs("args"))

	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
