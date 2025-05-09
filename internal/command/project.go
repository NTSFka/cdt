package command

import (
	"cdt/internal"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
)

var ProjectCommand = cli.Command{
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

func projectCommandTargetsAction(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	info, err := c.Project.Structure()

	if err != nil {
		return err
	}

	for name, target := range info.Targets {
		if target.Dependency {
			fmt.Printf("%v (dependency)\n", name)
		} else {
			fmt.Printf("%v\n", name)
		}
	}

	return nil
}

func projectCommandFilesAction(ctx context.Context, _ *cli.Command) error {
	c := ctx.Value("context").(internal.Context)
	info, err := c.Project.Structure()

	if err != nil {
		return err
	}

	for name, target := range info.Targets {
		if target.Dependency {
			continue
		}

		fmt.Printf("%v:\n", name)

		for _, file := range target.Files {
			fmt.Printf("  %v\n", file)
		}
	}

	return nil
}
