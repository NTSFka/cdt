package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"cdt/internal"
	"cdt/internal/workflow"
	"cdt/pkg"
)

// Release: `go build -ldflags="-X main.version=0.1.0"`.
var version = ""

func buildVersion() string {
	if version != "" {
		return version
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return "dev+" + setting.Value[0:7]
			}
		}
	}

	return "dev"
}

func parseEnvironment(environment string) (string, string) {
	parts := strings.SplitN(environment, ":", 2)

	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default: // nolint: gocritic
		fallthrough
	case 2:
		return parts[0], parts[1]
	}
}

func initEnvironment(
	directory string,
	environment *string,
	envProviders internal.EnvironmentProviders,
) (internal.Environment, error) {
	if environment != nil {
		toolName, argument := parseEnvironment(*environment)

		for _, provider := range envProviders {
			if provider.IsAvailable() &&
				(provider.Id() == toolName || slices.Contains(provider.Aliases(), toolName)) {
				return provider.CreateEnvironment(directory, argument)
			}
		}

		return nil, fmt.Errorf("environment '%s' not found", toolName)
	}

	// Try to detect environment automatically
	for _, provider := range envProviders {
		if env := provider.Detect(directory); env != nil {
			return *env, nil
		}
	}

	return internal.SystemEnvironment, nil
}

func buildContext(ctx context.Context, config internal.Config) (*internal.Context, error) {
	environmentProviders := pkg.InitEnvironmentProviders(
		ctx,
		internal.DetectOptions{
			ProjectDirectory: config.RootDirectory,
			Environment:      internal.SystemEnvironment,
			ToolsPaths:       config.Tools,
		},
	)

	env, err := initEnvironment(config.RootDirectory, config.Environment, environmentProviders)
	if err != nil {
		return nil, err
	}

	tools := pkg.InitTools(ctx, internal.DetectOptions{
		ProjectDirectory: config.RootDirectory,
		Environment:      env,
		ToolsPaths:       config.Tools,
	})

	project, err := workflow.CreateProject(config, tools)
	if err != nil {
		return nil, err
	}

	return &internal.Context{
		Config:               config,
		Project:              *project,
		Tools:                tools,
		EnvironmentProviders: environmentProviders,
		Environment:          env,
	}, nil
}

func main() {
	ctx := context.Background()

	app := pkg.NewApp(func(config internal.Config) (*internal.Context, error) {
		return buildContext(ctx, config)
	})

	app.Version = buildVersion()

	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
