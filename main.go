package main

import (
	"cdt/internal"
	"cdt/internal/tool"
	"cdt/internal/workflow"
	"cdt/pkg"
	"fmt"
	"os"
)

// InitTools initializes all supported tools on the system
func initTools(environment internal.Environment) internal.Tools {
	return internal.Tools{
		tool.DetectClangFormat(environment, nil),
		tool.DetectClangTidy(environment, nil),
		tool.DetectCMake(environment),
		tool.DetectCTest(environment),
	}
}

func detectProject(config internal.Config, tools internal.Tools) internal.Project {
	// CMake
	if p := workflow.DetectCMakeProject(config, tools); p != nil {
		return *p
	}

	return internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, internal.Workflow{})
}

func buildContext(config internal.Config) internal.Context {
	tools := initTools(internal.SystemEnvironment)

	return internal.Context{
		Config:  config,
		Project: detectProject(config, tools),
		Tools:   tools,
	}
}

func main() {
	if err := pkg.RunMain(buildContext, os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
