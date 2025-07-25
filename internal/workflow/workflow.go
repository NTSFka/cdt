package workflow

import (
	"cdt/internal"
	"errors"
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
	wf := internal.Workflow{}

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

// BuildProject creates a project from configuration and supported tools
func BuildProject(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	// User-defined workflow
	if config.Workflow != nil {
		switch cfg := config.Workflow.(type) {
		case *internal.ConfigWorkflow:
			if wf, err := FromConfig(*cfg, tools); err == nil {
				project := internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, *wf)
				return &project, nil
			} else {
				return nil, err
			}
		case string:
			// TODO: redesign
			switch cfg {
			case "cmake":
				if p := DetectCMakeProject(config, tools); p != nil {
					return p, nil
				} else {
					return nil, errors.New("cmake workflow requires CMakeLists.txt to be present in the project directory")
				}
			case "go":
				if p := DetectGoProject(config, tools); p != nil {
					return p, nil
				} else {
					return nil, errors.New("go workflow requires go.mod to be present in the project directory")
				}
			}

			return nil, fmt.Errorf("workflow '%s' not found", cfg)
		default:
			panic(fmt.Sprintf("unknown workflow type: %T", cfg))
		}
	}

	// CMake
	if p := DetectCMakeProject(config, tools); p != nil {
		return p, nil
	}

	// Go
	if p := DetectGoProject(config, tools); p != nil {
		return p, nil
	}

	project := internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, internal.Workflow{})

	return &project, nil
}
