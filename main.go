package main

import (
	"cdt/internal"
	"cdt/internal/project"
	"cdt/internal/tool"
	"cdt/pkg"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
)

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

func buildContext(config internal.Config) (*internal.Context, error) {
	environmentProviders := tool.InitEnvironmentProviders(internal.SystemEnvironment)

	env, err := initEnvironment(config.RootDirectory, config.Environment, environmentProviders)

	if err != nil {
		return nil, err
	}

	tools := tool.InitTools(env)

	p, err := project.BuildProject(config, tools)

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
	app := pkg.NewApp(buildContext)

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
