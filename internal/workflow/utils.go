package workflow

import (
	"cdt/internal"
	"fmt"
)

func getTool[T any](name *string, tools internal.Tools, what string) (*T, error) {
	if name == nil {
		return nil, nil
	}

	tool := tools.Get(*name)

	if tool == nil {
		return nil, fmt.Errorf("tool '%s' not found", *name)
	}

	if t, ok := tool.(T); ok {
		return &t, nil
	} else {
		return nil, fmt.Errorf("tool '%s' doesn't support %s", *name, what)
	}
}

// FromConfig creates workflow from configuration
//
//nolint:cyclop
func FromConfig(config internal.ConfigWorkflow, tools internal.Tools) (*internal.Workflow, error) {
	wf := internal.Workflow{
		Name: "custom",
	}

	if t, err := getTool[internal.ProjectConfigurator](config.Configure, tools, "configuration"); t != nil {
		wf.Configurator = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := getTool[internal.ProjectBuilder](config.Build, tools, "building"); t != nil {
		wf.Builder = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := getTool[internal.ProjectTester](config.Test, tools, "testing"); t != nil {
		wf.Tester = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := getTool[internal.ProjectFormatter](config.Format, tools, "formatting"); t != nil {
		wf.Formatter = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := getTool[internal.ProjectLinter](config.Lint, tools, "linting"); t != nil {
		wf.Linter = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := getTool[internal.ProjectRunner](config.Run, tools, "run"); t != nil {
		wf.Runner = *t
	} else if err != nil {
		return nil, err
	}

	return &wf, nil
}

func createProjectConfigCustom(config internal.Config, configWorkflow internal.ConfigWorkflow, tools internal.Tools) (*internal.Project, error) {
	if workflow, err := FromConfig(configWorkflow, tools); workflow != nil {
		return &internal.Project{
			Info: internal.ProjectInfo{
				Directory:         config.RootDirectory,
				StructureProvider: &internal.EmptyProjectStructureProvider{},
			},
			Workflow: *workflow,
		}, nil
	} else {
		return nil, err
	}
}

func createProjectConfigName(config internal.Config, workflowName string, tools internal.Tools) (*internal.Project, error) {
	for _, workflowType := range SupportedTypes {
		if workflowName == workflowType.Id() {
			cfg := Config{
				Directory:             config.RootDirectory,
				IntermediateDirectory: config.BuildDirectory,
			}

			workflow := workflowType.Create(cfg, tools)

			return &workflow, nil
		}
	}

	return nil, fmt.Errorf("workflow '%s' not found", workflowName)
}

func createProjectDetect(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	for _, pt := range SupportedTypes {
		if pt.Detect(config.RootDirectory) {
			cfg := Config{
				Directory:             config.RootDirectory,
				IntermediateDirectory: config.BuildDirectory,
			}

			workflow := pt.Create(cfg, tools)

			return &workflow, nil
		}
	}

	return &internal.Project{
		Info: internal.ProjectInfo{
			Directory:         config.RootDirectory,
			StructureProvider: &internal.EmptyProjectStructureProvider{},
		},
	}, nil
}

// CreateProject creates a project from configuration and supported tools
func CreateProject(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	// User-defined workflow
	if config.Workflow != nil {
		switch cfg := config.Workflow.(type) {
		case *internal.ConfigWorkflow:
			return createProjectConfigCustom(config, *cfg, tools)
		case string:
			return createProjectConfigName(config, cfg, tools)
		default:
			panic(fmt.Sprintf("unknown workflow type: %T", cfg))
		}
	}

	return createProjectDetect(config, tools)
}
