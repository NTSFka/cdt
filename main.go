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
func initTools() Tools {
	return Tools{
		tool.DetectClangFormat(nil),
		tool.DetectClangTidy(nil),
		tool.DetectCMake(),
		tool.DetectCTest(),
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
	tools := initTools()

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
