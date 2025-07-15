package main

import (
	"cdt/internal"
	"cdt/internal/tool"
	"cdt/internal/workflow"
	"cdt/pkg"
	"errors"
	"fmt"
	"os"
	"strings"
)

var environmentProviders = internal.EnvironmentProviders{
	tool.DetectDocker(internal.SystemEnvironment),
	tool.DetectDockerCompose(internal.SystemEnvironment),
}

func parseEnvironment(environment string) (string, string, error) {
	parts := strings.SplitN(environment, ":", 2)

	if len(parts) != 2 {
		return "", "", errors.New("invalid environment string")
	}

	return parts[0], parts[1], nil
}

func initEnvironment(rootDirectory string, environment *string) (internal.Environment, error) {
	if environment != nil {
		toolName, argument, err := parseEnvironment(*environment)

		if err != nil {
			return nil, err
		}

		for _, provider := range environmentProviders {
			if provider.IsAvailable() && provider.Id() == toolName {
				return provider.CreateEnvironment(rootDirectory, argument)
			}
		}
	}

	return internal.SystemEnvironment, nil
}

func detectProject(config internal.Config, tools tool.SupportedTools) internal.Project {
	// CMake
	if p := workflow.DetectCMakeProject(config, tools); p != nil {
		return *p
	}

	// Go
	if p := workflow.DetectGoProject(config, tools); p != nil {
		return *p
	}

	return internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, internal.Workflow{})
}

func buildContext(config internal.Config) (*internal.Context, error) {
	env, err := initEnvironment(config.RootDirectory, config.Environment)

	if err != nil {
		return nil, err
	}

	tools := tool.NewSupportedTools(env)

	return &internal.Context{
		Config:               config,
		Project:              detectProject(config, tools),
		Tools:                tools.ToTools(),
		EnvironmentProviders: environmentProviders,
		Environment:          env,
	}, nil
}

func main() {
	if err := pkg.RunMain(buildContext, os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
