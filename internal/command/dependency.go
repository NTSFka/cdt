package command

import (
	"cdt/internal"
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewDependencyCommand() *cli.Command {
	return &cli.Command{
		Name:    "dependency",
		Usage:   "manage project dependencies",
		Aliases: []string{"dep", "d"},
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "add dependencies to the project",
				Action: dependencyAddCommandAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dev",
						Usage: "add dependencies to the project development dependencies",
						Value: false,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArgs{
						Name: "dependencies",
						Min:  1,
						Max:  -1,
					},
				},
			},
			{
				Name:   "remove",
				Usage:  "remove dependencies from the project",
				Action: dependencyRemoveCommandAction,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "dev",
						Usage: "remove dependencies from the project development dependencies",
						Value: false,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArgs{
						Name: "dependencies",
						Min:  1,
						Max:  -1,
					},
				},
			},
			{
				Name:   "update",
				Usage:  "update dependencies in the project",
				Action: dependencyUpdateCommandAction,
				Arguments: []cli.Argument{
					&cli.StringArgs{
						Name: "dependencies",
						Min:  0,
						Max:  -1,
					},
				},
			},
			{
				Name:  "fetch",
				Usage: "fetch the project dependencies",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "no-dev",
						Usage: "fetch the project dependencies without development dependencies",
						Value: false,
					},
				},
				Action: dependencyFetchCommandAction,
			},
			{
				Name:   "list",
				Usage:  "list the project dependencies",
				Action: dependencyListCommandAction,
			},
			{
				Name:   "audit",
				Usage:  "audit the project dependencies",
				Action: dependencyAuditCommandAction,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific dependency tool",
			},
		},
	}
}

func dependencyAddCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	manager, err := getDependencyManager(cmdContext, cmd)

	if err != nil {
		return err
	}

	dependencies := cmd.StringArgs("dependencies")
	dev := cmd.Bool("dev")

	options := internal.ProjectDependencyManagerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	if err := manager.AddDependencies(ctx, options, dependencies, dev); err != nil {
		return fmt.Errorf("failed to add dependencies: %w", err)
	}

	return nil
}

func dependencyRemoveCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	manager, err := getDependencyManager(cmdContext, cmd)

	if err != nil {
		return err
	}

	dependencies := cmd.StringArgs("dependencies")
	dev := cmd.Bool("dev")

	options := internal.ProjectDependencyManagerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	if err := manager.RemoveDependencies(ctx, options, dependencies, dev); err != nil {
		return fmt.Errorf("failed to remove dependencies: %w", err)
	}

	return nil
}

func dependencyUpdateCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	manager, err := getDependencyManager(cmdContext, cmd)

	if err != nil {
		return err
	}

	dependencies := cmd.StringArgs("dependencies")

	options := internal.ProjectDependencyManagerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	if err := manager.UpdateDependencies(ctx, options, dependencies); err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}

	return nil
}

func dependencyFetchCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	manager, err := getDependencyManager(cmdContext, cmd)

	if err != nil {
		return err
	}

	noDev := cmd.Bool("no-dev")

	options := internal.ProjectDependencyManagerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	if err := manager.FetchDependencies(ctx, options, noDev); err != nil {
		return fmt.Errorf("failed to fetch dependencies: %w", err)
	}

	return nil
}

func dependencyListCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	manager, err := getDependencyManager(cmdContext, cmd)

	if err != nil {
		return err
	}

	options := internal.ProjectDependencyManagerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	if err := manager.ListDependencies(ctx, options); err != nil {
		return fmt.Errorf("failed to list dependencies: %w", err)
	}

	return nil
}

func dependencyAuditCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	manager, err := getDependencyManager(cmdContext, cmd)

	if err != nil {
		return err
	}

	options := internal.ProjectDependencyManagerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
	}

	if err := manager.AuditDependencies(ctx, options); err != nil {
		return fmt.Errorf("failed to audit dependencies: %w", err)
	}

	return nil
}

func getDependencyManager(c internal.Context, cmd *cli.Command) (internal.ProjectDependencyManager, error) {
	manager := c.Project.Workflow.DependencyManager

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := c.Tools.Get(toolName)

		if tool == nil {
			return nil, fmt.Errorf("tool '%s' not found", toolName)
		}

		managerTool, ok := tool.(internal.ProjectDependencyManager)

		if ok {
			manager = managerTool
		} else {
			return nil, fmt.Errorf("tool '%s' doesn't support dependency management", toolName)
		}
	}

	if manager == nil {
		return nil, errors.New("project doesn't support dependency management")
	}

	return manager, nil
}
