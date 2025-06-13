package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

var LintCommand = cli.Command{
	Name:   "lint",
	Usage:  "Lint the project",
	Action: lintCommandAction,
	Arguments: []cli.Argument{
		&cli.StringArgs{
			Name: "files",
			Min:  0,
			Max:  -1,
		},
	},
}

func lintCommandAction(ctx context.Context, command *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	linter := c.Project.Workflow.Linter

	if linter == nil {
		return errors.New("project doesn't support linting")
	}

	var err error
	if files := command.StringArgs("files"); len(files) > 0 {
		err = linter.LintFiles(c.Project, files, command.Args().Tail())
	} else {
		err = linter.LintAll(c.Project, command.Args().Tail())
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
