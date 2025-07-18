package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

var ExecCommand = cli.Command{
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

func execCommandAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)

	command := cmd.StringArg("command")

	if command == "" {
		return errors.New("COMMAND is required")
	}

	options := internal.RunOptions{
		Directory: c.Project.RootDirectory(),
		Input:     cmd.Reader,
		Output:    cmd.Writer,
		Error:     cmd.ErrWriter,
	}

	err := c.Environment.RunExecutable(ctx, options, command, cmd.StringArgs("args"))

	if err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
