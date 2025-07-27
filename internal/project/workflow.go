package project

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
