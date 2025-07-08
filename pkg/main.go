package pkg

import (
	"cdt/internal"
	"cdt/internal/command"
	"context"
	"errors"
	"github.com/urfave/cli/v3"
	"strings"
)

func parseEnvironment(environment string) (*internal.ConfigEnvironment, error) {
	parts := strings.SplitN(environment, ":", 2)

	if len(parts) != 2 {
		return nil, errors.New("invalid environment string")
	}

	return &internal.ConfigEnvironment{
		ToolName: parts[0],
		Argument: parts[1],
	}, nil
}

// RunMain is the main function of the application
func RunMain(buildContext func(config internal.Config) (*internal.Context, error), args []string) error {
	cmd := &cli.Command{
		Name:                  "cdt",
		Usage:                 "A common developer tool",
		Version:               "0.1.0",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"d"},
				Usage:   "project directory",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "build directory",
				Value:   "build",
			},
			&cli.StringFlag{
				Name:    "environment",
				Aliases: []string{"e"},
				Usage:   "environment to use",
				Value:   "",
			},
		},
		Commands: []*cli.Command{
			&command.ProjectCommand,
			&command.ToolCommand,
			&command.ConfigureCommand,
			&command.BuildCommand,
			&command.FormatCommand,
			&command.TestCommand,
			&command.LintCommand,
			&command.RunCommand,
			&command.EnvironmentCommand,
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			projectPath := cmd.String("directory")
			var buildDirectory *string = nil

			if cmd.Count("build") > 0 {
				directory := cmd.String("build")
				buildDirectory = &directory
			}

			var environment *internal.ConfigEnvironment

			if cmd.Count("environment") > 0 {
				env, err := parseEnvironment(cmd.String("environment"))

				if err != nil {
					return nil, err
				}

				environment = env
			}

			config := internal.Config{
				RootDirectory:  projectPath,
				BuildDirectory: buildDirectory,
				Environment:    environment,
			}

			c, err := buildContext(config)

			if err != nil {
				return nil, err
			}

			return context.WithValue(ctx, "context", *c), nil //nolint:staticcheck
		},
		After: func(ctx context.Context, command *cli.Command) error {
			c, ok := ctx.Value("context").(internal.Context)

			if ok && c.Environment != nil {
				return c.Environment.Cleanup(ctx)
			}

			return nil
		},
	}

	return cmd.Run(context.Background(), args)
}
