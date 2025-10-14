package command_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/command"
	"cdt/internal/test"
	"cdt/internal/tool"
	"cdt/internal/workflow"

	"github.com/stretchr/testify/require"
)

func runWorkflow(ctx context.Context, context internal.Context, args ...string) error {
	return test.RunCommand(ctx, command.NewWorkflowCommand(), context, args...)
}

func TestWorkflowList(t *testing.T) {
	err := runWorkflow(t.Context(), internal.Context{}, "list")

	require.NoError(t, err)
}

func TestWorkflowShow_Empty(t *testing.T) {
	err := runWorkflow(t.Context(), internal.Context{}, "show")

	require.NoError(t, err)
}

func TestWorkflowShow_Custom(t *testing.T) {
	configurator := struct {
		internal.ExecutableTool
		test.ProjectConfigurator
	}{
		internal.MakeExecutableTool(
			"test1",
			"Test",
			"",
			internal.Tags{},
			test.LazyExecutableNil,
		),
		test.ProjectConfigurator{},
	}

	linter1 := struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool(
			"test1",
			"Test 1",
			"",
			internal.Tags{},
			test.LazyExecutableNil,
		),
		test.ProjectLinter{},
	}

	linter2 := struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool(
			"test2",
			"Test 2",
			"",
			internal.Tags{},
			test.LazyExecutableNil,
		),
		test.ProjectLinter{},
	}

	err := runWorkflow(t.Context(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Name:              "custom",
				Configurator:      &configurator,
				Linter:            &workflow.LinterList{&linter1, &linter2},
				DependencyManager: test.NewDependencyManager(t), // unknown tool
			},
		},
	}, "show")

	require.NoError(t, err)
}

func TestWorkflowShow_Named(t *testing.T) {
	tools := internal.Tools{
		tool.NewGo(test.LazyExecutable("go-test")),
		tool.NewGolangCILint(test.LazyExecutableNil),
	}

	err := runWorkflow(t.Context(), internal.Context{
		Tools: tools,
	}, "show", "go")

	require.NoError(t, err)
}
