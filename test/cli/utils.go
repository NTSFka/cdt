package cli

import (
	"cdt/internal"
	"cdt/pkg"
	"errors"
)

func runMain(contextBuilder func(config internal.Config) (*internal.Context, error), args ...string) error {
	return pkg.RunMain(contextBuilder, append([]string{"cdt"}, args...))
}

// Run main function and return obtained configuration
func runMainGetConfig(args ...string) (config internal.Config) {
	_ = runMain(func(cfg internal.Config) (*internal.Context, error) {
		config = cfg
		return nil, errors.New("not used")
	}, args...)

	return
}

func runMainWithEnvironment(environment internal.Environment, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:      internal.Config{},
			Project:     internal.Project{},
			Tools:       nil,
			Environment: environment,
		}, nil
	}, args...)
}

func runMainWithProject(project internal.Project, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:  internal.Config{},
			Project: project,
			Tools:   nil,
		}, nil
	}, args...)
}

func runMainWithWorkflow(workflow internal.Workflow, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:  internal.Config{},
			Project: internal.Project{Workflow: workflow},
			Tools:   nil,
		}, nil
	}, args...)
}

func runMainWithTools(tools internal.Tools, args ...string) error {
	return runMain(func(config internal.Config) (*internal.Context, error) {
		return &internal.Context{
			Config:  internal.Config{},
			Project: internal.Project{},
			Tools:   tools,
		}, nil
	}, args...)
}
