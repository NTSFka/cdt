package main

import (
	"cdt/internal"
	"cdt/internal/tool"
	"cdt/internal/workflow"
	"cdt/pkg"
	"fmt"
	"os"
)

var environmentProviders = internal.EnvironmentProviders{
	tool.DetectDocker(),
}

func initEnvironment(rootDirectory string, environment *internal.ConfigEnvironment) (internal.Environment, error) {
	if environment != nil {
		for _, provider := range environmentProviders {
			if provider.IsAvailable() && provider.Id() == environment.ToolName {
				return provider.CreateEnvironment(rootDirectory, environment.Argument)
			}
		}
	}

	return internal.SystemEnvironment, nil
}

// InitTools initializes all supported tools on the system
func initTools(environment internal.Environment) tool.SupportedTools {
	return tool.SupportedTools{
		ClangFormat:  tool.DetectClangFormat(environment, nil),
		ClangTidy:    tool.DetectClangTidy(environment, nil),
		CMake:        tool.DetectCMake(environment),
		CTest:        tool.DetectCTest(environment),
		Go:           tool.DetectGo(environment),
		GolangCILint: tool.DetectGolangCILint(environment),
	}
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

	tools := initTools(env)

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
