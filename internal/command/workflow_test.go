package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"cdt/internal/workflow"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runWorkflow(ctx context.Context, context internal.Context, args ...string) error {
	return test.RunCommand(ctx, NewWorkflowCommand(), context, args...)
}

func TestWorkflowList(t *testing.T) {
	err := runWorkflow(context.Background(), internal.Context{}, "list")

	assert.NoError(t, err)
}

func TestWorkflowShow_Empty(t *testing.T) {
	err := runWorkflow(context.Background(), internal.Context{}, "show")

	assert.NoError(t, err)
}

func TestWorkflowShow_Custom(t *testing.T) {
	configurator := struct {
		internal.ExecutableTool
		test.ProjectConfigurator
	}{
		internal.MakeExecutableTool("test1", "Test", "", internal.Tags{}, func() *internal.Executable {
			return nil
		}),
		test.ProjectConfigurator{},
	}

	linter1 := struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool("test1", "Test 1", "", internal.Tags{}, func() *internal.Executable {
			return nil
		}),
		test.ProjectLinter{},
	}

	linter2 := struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool("test2", "Test 2", "", internal.Tags{}, func() *internal.Executable {
			return nil
		}),
		test.ProjectLinter{},
	}

	err := runWorkflow(context.Background(), internal.Context{
		Project: internal.Project{
			Workflow: internal.Workflow{
				Name:              "custom",
				Configurator:      &configurator,
				Linter:            &workflow.LinterList{&linter1, &linter2},
				DependencyManager: test.NewDependencyManager(t), // unknown tool
			},
		},
	}, "show")

	assert.NoError(t, err)
}

func TestWorkflowShow_Named(t *testing.T) {
	tools := internal.Tools{
		tool.NewGo(func() *internal.Executable { return &internal.Executable{Path: "go-test"} }),
		tool.NewGolangCILint(func() *internal.Executable { return nil }),
	}

	err := runWorkflow(context.Background(), internal.Context{
		Tools: tools,
	}, "show", "go")

	assert.NoError(t, err)
}
