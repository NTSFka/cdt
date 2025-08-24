package pkg

import (
	"cdt/internal"
	"cdt/internal/command"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

const ConfigFileName = "cdt.yml"

// NewApp creates the main CLI application
func NewApp(buildContext func(config internal.Config) (*internal.Context, error)) *cli.Command {
	return &cli.Command{
		Name:                  "cdt",
		Usage:                 "A common developer tool",
		Version:               "0.1.0",
		EnableShellCompletion: true,
		Authors:               []any{"Jiří Fatka"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "root",
				Aliases: []string{"r"},
				Usage:   "Project root directory. If contains " + ConfigFileName + " file, it will be used as a configuration file.",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Build directory for intermediate data. If not specified the value from configuration file some default will be used.",
				Value:   "build",
			},
			&cli.StringFlag{
				Name:    "environment",
				Aliases: []string{"e"},
				Usage:   "Environment to use, e.g. `docker:image`. If not specified the value from configuration file or system environment will be used.",
				Value:   "system",
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "What project type to use, e.g. `go`. If not specified the value from configuration file or detected type will be used.",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to configuration file. By default `cdt.yml` in root directory will be used, if present.",
				Value:   "cdt.yml",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Enable debug output",
			},
		},
		Commands: []*cli.Command{
			command.NewProjectCommand(),
			command.NewToolCommand(),
			command.NewConfigureCommand(),
			command.NewBuildCommand(),
			command.NewFormatCommand(),
			command.NewTestCommand(),
			command.NewLintCommand(),
			command.NewRunCommand(),
			command.NewEnvironmentCommand(),
			command.NewExecCommand(),
			command.NewDependencyCommand(),
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
}

func createConfig(cmd *cli.Command) (*internal.Config, error) {
	config := internal.DefaultConfig()
	var configPath string

	// The root directory is not affected by the configuration file (it's used as a search path)
	if cmd.Count("root") > 0 {
		config.RootDirectory = cmd.String("root")
	}

	if cmd.Count("config") > 0 {
		configPath = cmd.String("config")
	} else {
		configPath = filepath.Join(config.RootDirectory, ConfigFileName)
	}

	fileConfig, err := loadConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	if fileConfig != nil {
		fileConfig.UpdateConfig(&config)
	}

	// Override configuration file
	if cmd.Count("type") > 0 {
		config.Workflow = cmd.String("type")
	}

	// Override configuration file
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
