package main

import (
	"cdt/internal"
	"cdt/internal/project"
	"cdt/internal/tool"
	"cdt/pkg"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

func parseEnvironment(environment string) (string, string, error) {
	parts := strings.SplitN(environment, ":", 2)

	if len(parts) != 2 {
		return "", "", errors.New("invalid environment string")
	}

	return parts[0], parts[1], nil
}

func initEnvironment(rootDirectory string, environment *string, environmentProviders internal.EnvironmentProviders) (internal.Environment, error) {
	if environment != nil {
		toolName, argument, err := parseEnvironment(*environment)

		if err != nil {
			return nil, err
		}

		for _, provider := range environmentProviders {
			if provider.IsAvailable() && (provider.Id() == toolName || provider.IdShort() == toolName) {
				return provider.CreateEnvironment(rootDirectory, argument)
			}
		}

		return nil, fmt.Errorf("environment '%s' not found", toolName)
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
