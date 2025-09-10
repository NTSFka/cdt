package command

import (
	"cdt/internal"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

func NewProjectCommand() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "Show project information",
		Commands: []*cli.Command{
			{
				Name:   "targets",
				Usage:  "show available targets",
				Action: projectCommandTargetsAction,
			},
			{
				Name:   "files",
				Usage:  "show project files",
				Action: projectCommandFilesAction,
			},
		},
	}
}

func projectCommandTargetsAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	info, err := c.ProjectInfo.Structure()

	if err != nil {
		return err
	}

	for name, target := range info.Targets {
		if target.Dependency {
			_, _ = fmt.Fprintf(cmd.Writer, "%v (dependency)\n", name)
		} else {
			_, _ = fmt.Fprintf(cmd.Writer, "%v\n", name)
		}
	}

	return nil
}

func projectCommandFilesAction(ctx context.Context, cmd *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	info, err := c.ProjectInfo.Structure()

	if err != nil {
		return err
	}

	for name, target := range info.Targets {
		if target.Dependency {
			continue
		}

		_, _ = fmt.Fprintf(cmd.Writer, "%v:\n", name)

		for _, file := range target.Files {
			_, _ = fmt.Fprintf(cmd.Writer, "  %v\n", file)
		}
	}

	return nil
}
