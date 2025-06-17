package main

import (
	. "cdt/internal"
	"cdt/internal/project"
	"cdt/internal/tool"
	. "cdt/pkg"
	"fmt"
	"os"
)

// InitTools initializes all supported tools on the system
func initTools(environment Environment) Tools {
	return Tools{
		tool.DetectClangFormat(environment, nil),
		tool.DetectClangTidy(environment, nil),
		tool.DetectCMake(environment),
		tool.DetectCTest(environment),
	}
}

func detectProject(config Config, tools Tools) Project {
	// CMake
	if p := project.DetectCMakeProject(config, tools); p != nil {
		return *p
	}

	return MakeProject(config.RootDirectory, "", &EmptyProjectStructureProvider{}, Workflow{})
}

func buildContext(config Config) Context {
	tools := initTools(SystemEnvironment)

	return Context{
		Config:  config,
		Project: detectProject(config, tools),
		Tools:   tools,
	}
}

func main() {
	if err := RunMain(buildContext, os.Args); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}
