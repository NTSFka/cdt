package pkg

import (
	"cdt/internal"
	"cdt/internal/command"
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
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
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "configuration file",
				Value:   "cdt.yml",
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

func determineProjectDir(cmd *cli.Command, config *internal.FileConfig) (result string) {
	result = "."

	if config != nil && config.Project.Directory != nil {
		result = *config.Project.Directory
	}

	if cmd.Count("directory") > 0 {
		result = cmd.String("directory")
	}

	return
}

func determineBuildDir(cmd *cli.Command, config *internal.FileConfig) (result *string) {
	result = nil

	if config != nil && config.Project.BuildDirectory != nil {
		result = config.Project.BuildDirectory
	}

	if cmd.Count("build") > 0 {
		directory := cmd.String("build")
		result = &directory
	}

	return
}

func determineEnvironment(cmd *cli.Command, config *internal.FileConfig) (*internal.ConfigEnvironment, error) {
	if cmd.Count("environment") > 0 {
		return parseEnvironment(cmd.String("environment"))
	}

	if config != nil && config.Project.Environment != nil {
		return parseEnvironment(*config.Project.Environment)
	}

	return nil, nil
}

func createConfig(cmd *cli.Command) (*internal.Config, error) {
	fileConfig, err := loadConfigFile(cmd.String("config"))
	if err != nil {
		return nil, err
	}

	projectPath := determineProjectDir(cmd, fileConfig)
	buildDirectory := determineBuildDir(cmd, fileConfig)
	environment, err := determineEnvironment(cmd, fileConfig)

	if err != nil {
		return nil, err
	}

	return &internal.Config{
		RootDirectory:  projectPath,
		BuildDirectory: buildDirectory,
		Environment:    environment,
	}, nil
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
