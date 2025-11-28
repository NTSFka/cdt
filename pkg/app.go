package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"cdt/internal"
	"cdt/internal/command"

	"github.com/urfave/cli/v3"
)

const ConfigFileName = "cdt.yml"

// NewApp creates the main CLI application.
// nolint:funlen
func NewApp(buildContext func(config internal.Config) (*internal.Context, error)) *cli.Command {
	return &cli.Command{
		Name:                  "cdt",
		Usage:                 "A common developer tool",
		EnableShellCompletion: true,
		Authors:               []any{"Jiří Fatka"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "root",
				Aliases: []string{"r"},
				Usage: "Project root directory. If contains " + ConfigFileName + " file, it will be used as " +
					"a configuration file.",
				Value: ".",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage: "Build directory for output data. If not specified the value from configuration file " +
					"some default might be used.",
			},
			&cli.StringFlag{
				Name:    "environment",
				Aliases: []string{"e"},
				Usage: "Environment to use, e.g. `docker:image`. If not specified the value from configuration file " +
					"or system environment will be used.",
				Value: "system",
			},
			&cli.StringFlag{
				Name:    "workflow",
				Aliases: []string{"w"},
				Usage: "What workflow to use, e.g. `go`. If not specified the value from configuration file or " +
					"detected workflow will be used.",
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
			command.NewWorkflowCommand(),
			command.NewCoverageCommand(),
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			if cmd.Bool("debug") {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}

			config, err := createConfig(cmd)

			if err != nil {
				return nil, err
			}

			cmdContext, err := buildContext(*config)

			if err != nil {
				return nil, err
			}

			return context.WithValue(ctx, "context", *cmdContext), nil //nolint:staticcheck
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
	if cmd.Count("workflow") > 0 {
		config.Workflow = cmd.String("workflow")
	}

	// Override configuration file
	if cmd.Count("output") > 0 {
		directory := cmd.String("output")
		config.OutputDirectory = &directory
	}

	if cmd.Count("environment") > 0 {
		environment := cmd.String("environment")
		config.Environment = &environment
	}

	return &config, nil
}

func loadConfigFile(configPath string) (*internal.FileConfig, error) {
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // nolint: nilnil
		}

		return nil, fmt.Errorf("failed to open configuration file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	fileConfig, err := internal.LoadConfigFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration file: %w", err)
	}

	return fileConfig, nil
}
