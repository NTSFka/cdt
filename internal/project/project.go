package project

import (
	"cdt/internal"
	"fmt"
)

var projectTypes = map[string]func(internal.Config, internal.Tools) (*internal.Project, error){
	"go":    DetectGoProject,
	"cmake": DetectCMakeProject,
}

func buildProjectConfigCustom(config internal.Config, cwf internal.ConfigWorkflow, tools internal.Tools) (*internal.Project, error) {
	if wf, err := FromConfig(cwf, tools); wf != nil {
		project := internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, *wf)
		return &project, nil
	} else {
		return nil, err
	}
}

func buildProjectConfigName(config internal.Config, workflow string, tools internal.Tools) (*internal.Project, error) {
	for name, f := range projectTypes {
		if workflow == name {
			if p, err := f(config, tools); p != nil {
				return p, nil
			} else {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("workflow '%s' not found", workflow)
}

func buildProjectDetect(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	for _, f := range projectTypes {
		if p, _ := f(config, tools); p != nil {
			return p, nil
		}
	}

	project := internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, internal.Workflow{})

	return &project, nil
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
