package pkg

import (
	"cdt/internal"
	"cdt/internal/command"
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log/slog"
	"os"
)

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
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "configuration file",
				Value:   "cdt.yml",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug output",
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
			&command.ExecCommand,
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.Bool("debug") {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}

			config, err := createConfig(cmd)

			if err != nil {
				return nil, err
			}

			c, err := buildContext(*config)

			if err != nil {
				return nil, err
			}

			return context.WithValue(ctx, "context", *c), nil //nolint:staticcheck
		},
		After: func(ctx context.Context, cmd *cli.Command) error {
			c, ok := ctx.Value("context").(internal.Context)

			if ok && c.Environment != nil {
				return c.Environment.Cleanup(ctx)
			}

			return nil
		},
	}

	return cmd.Run(context.Background(), args)
}

func createConfig(cmd *cli.Command) (*internal.Config, error) {
	config := internal.DefaultConfig()

	fileConfig, err := loadConfigFile(cmd.String("config"))
	if err != nil {
		return nil, err
	}

	if fileConfig != nil {
		fileConfig.UpdateConfig(&config)
	}

	if cmd.Count("directory") > 0 {
		config.RootDirectory = cmd.String("directory")
	}

	if cmd.Count("build") > 0 {
		directory := cmd.String("build")
		config.BuildDirectory = &directory
	}

	if cmd.Count("environment") > 0 {
		environment := cmd.String("environment")
		config.Environment = &environment
	}

	return &config, nil
}

func loadConfigFile(configPath string) (*internal.FileConfig, error) {
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to open configuration file: %w", err)
	}

	fileConfig, err := internal.LoadConfigFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration file: %w", err)
	}

	return fileConfig, nil
}
