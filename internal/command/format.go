package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
)

var FormatCommand = cli.Command{
	Name:   "format",
	Usage:  "Format files in the project",
	Action: formatCommandAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "check",
			Value: false,
			Usage: "Check if the project or given files are formatted",
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

func formatCommandAction(ctx context.Context, command *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	formatter := c.Project.Formatter()

	if formatter == nil {
		return errors.New("project doesn't support source formatting")
	}

	var err error
	if command.Bool("check") {
		if files := command.StringArgs("files"); len(files) > 0 {
			err = formatter.FormatCheckFiles(files)
		} else {
			err = formatter.FormatCheckAll()
		}
	} else {
		if files := command.StringArgs("files"); len(files) > 0 {
			err = formatter.FormatFiles(files)
		} else {
			err = formatter.FormatAll()
		}
	}

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
