package main

import (
	"cdt/internal"
	"cdt/internal/tool"
	"cdt/internal/workflow"
	"cdt/pkg"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
)

var version = "dev"

func parseEnvironment(environment string) (string, string) {
	parts := strings.SplitN(environment, ":", 2)

	if len(parts) == 1 {
		return parts[0], ""
	}

	return parts[0], parts[1]
}

func initEnvironment(directory string, environment *string, envProviders internal.EnvironmentProviders) (internal.Environment, error) {
	if environment != nil {
		toolName, argument := parseEnvironment(*environment)

		for _, provider := range envProviders {
			if provider.IsAvailable() && (provider.Id() == toolName || slices.Contains(provider.Aliases(), toolName)) {
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
	environmentProviders := tool.InitEnvironmentProviders(ctx, internal.SystemEnvironment)

	env, err := initEnvironment(config.RootDirectory, config.Environment, environmentProviders)

	if err != nil {
		return nil, err
	}

	tools := tool.InitTools(ctx, env)

	p, err := workflow.CreateProject(config, tools)

	if err != nil {
		return nil, err
	}

	return &internal.Context{
		Config:               config,
		Project:              *p,
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

	app.Version = version

	if err := app.Run(ctx, os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
