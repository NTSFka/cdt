package main

import (
	"cdt/internal"
	"cdt/internal/tool"
	"cdt/internal/workflow"
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
			if provider.IsAvailable() && provider.Id() == toolName {
				return provider.CreateEnvironment(rootDirectory, argument)
			}
		}
	}

	return internal.SystemEnvironment, nil
}

func buildWorkflowGetTool[T any](name *string, tools internal.Tools, what string) (*T, error) {
	if name != nil {
		if t, ok := tools.Get(*name).(T); ok {
			return &t, nil
		} else {
			return nil, fmt.Errorf("tool '%s' doesn't support %s", *name, what)
		}
	}

	return nil, nil
}

//nolint:cyclop
func buildWorkflow(config internal.ConfigWorkflow, tools internal.Tools) (*internal.Workflow, error) {
	wf := internal.Workflow{}

	if t, err := buildWorkflowGetTool[internal.ProjectConfigurator](config.Configure, tools, "configuration"); t != nil {
		wf.Configurator = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := buildWorkflowGetTool[internal.ProjectBuilder](config.Build, tools, "building"); t != nil {
		wf.Builder = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := buildWorkflowGetTool[internal.ProjectTester](config.Test, tools, "testing"); t != nil {
		wf.Tester = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := buildWorkflowGetTool[internal.ProjectFormatter](config.Format, tools, "formatting"); t != nil {
		wf.Formatter = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := buildWorkflowGetTool[internal.ProjectLinter](config.Lint, tools, "linting"); t != nil {
		wf.Linter = *t
	} else if err != nil {
		return nil, err
	}

	if t, err := buildWorkflowGetTool[internal.ProjectRunner](config.Run, tools, "run"); t != nil {
		wf.Runner = *t
	} else if err != nil {
		return nil, err
	}

	return &wf, nil
}

func detectProject(config internal.Config, tools internal.Tools) (*internal.Project, error) {
	if cwf := config.Workflow; cwf != nil {
		if wf, err := buildWorkflow(*cwf, tools); err != nil {
			return nil, err
		} else {
			project := internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, *wf)
			return &project, nil
		}
	}

	// CMake
	if p := workflow.DetectCMakeProject(config, tools); p != nil {
		return p, nil
	}

	// Go
	if p := workflow.DetectGoProject(config, tools); p != nil {
		return p, nil
	}

	project := internal.MakeProject(config.RootDirectory, "", &internal.EmptyProjectStructureProvider{}, internal.Workflow{})

	return &project, nil
}

func buildContext(config internal.Config) (*internal.Context, error) {
	environmentProviders := tool.InitEnvironmentProviders(internal.SystemEnvironment)

	env, err := initEnvironment(config.RootDirectory, config.Environment, environmentProviders)

	if err != nil {
		return nil, err
	}

	tools := tool.InitTools(env)

	project, err := detectProject(config, tools)

	if err != nil {
		return nil, err
	}

	return &internal.Context{
		Config:               config,
		Project:              *project,
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
