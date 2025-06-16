package pkg

import (
	. "cdt/internal"
	"cdt/internal/command"
	"context"
	"github.com/urfave/cli/v3"
)

// RunMain is the main function of the application
func RunMain(buildContext func(config Config) Context, args []string) error {
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
		},
		Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
			projectPath := command.String("directory")
			var buildDirectory *string = nil

			if command.Count("build") > 0 {
				directory := command.String("build")
				buildDirectory = &directory
			}

			config := Config{
				RootDirectory:  projectPath,
				BuildDirectory: buildDirectory,
			}

			c := buildContext(config)

			return context.WithValue(ctx, "context", c), nil
		},
	}

	return cmd.Run(context.Background(), args)
}
