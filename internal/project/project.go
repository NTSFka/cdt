package project

import (
	"cdt/internal"
	"fmt"
)

var projectTypes = []Type{
	&GoType{},
	&CMakeType{},
	&PHPType{},
	&PythonType{},
}

type Config struct {
	Directory      string
	BuildDirectory *string
}

// Type specifies a predefined project type
type Type interface {
	// Id returns project type unique identifier.
	Id() string

	// Detect detects if a project of this type is in a given directory.
	Detect(directory string) bool

	// Create a project from a directory with a given tool set.
	Create(config Config, tools internal.Tools) internal.Project
}

func buildProjectConfigCustom(config internal.Config, cwf internal.ConfigWorkflow, tools internal.Tools) (*internal.Project, error) {
	if wf, err := FromConfig(cwf, tools); wf != nil {
		return &internal.Project{
			Type: "custom",
			Desc: internal.ProjectInfo{
				Directory:         config.RootDirectory,
				StructureProvider: &internal.EmptyProjectStructureProvider{},
			},
			Workflow: *wf,
		}, nil
	} else {
		return nil, err
	}
}

func buildProjectConfigName(config internal.Config, workflow string, tools internal.Tools) (*internal.Project, error) {
	for _, pt := range projectTypes {
		if workflow == pt.Id() {
			project := pt.Create(Config{Directory: config.RootDirectory, BuildDirectory: config.BuildDirectory}, tools)
			return &project, nil
		}
	}

	return nil, fmt.Errorf("workflow '%s' not found", workflow)
}

func buildProjectDetect(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	for _, pt := range projectTypes {
		if pt.Detect(config.RootDirectory) {
			project := pt.Create(Config{Directory: config.RootDirectory, BuildDirectory: config.BuildDirectory}, tools)
			return &project, nil
		}
	}

	return &internal.Project{
		Type: "unknown",
		Desc: internal.ProjectInfo{
			Directory:         config.RootDirectory,
			StructureProvider: &internal.EmptyProjectStructureProvider{},
		},
	}, nil
}

// BuildProject creates a project from configuration and supported tools
func BuildProject(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	// User-defined workflow
	if config.Workflow != nil {
		switch cfg := config.Workflow.(type) {
		case *internal.ConfigWorkflow:
			return buildProjectConfigCustom(config, *cfg, tools)
		case string:
			return buildProjectConfigName(config, cfg, tools)
		default:
			panic(fmt.Sprintf("unknown workflow type: %T", cfg))
		}
	}

	return buildProjectDetect(config, tools)
}
