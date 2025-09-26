package command

import (
	"context"
	"fmt"
	"io"

	"cdt/internal"
	"cdt/internal/workflow"

	"github.com/urfave/cli/v3"
)

func NewWorkflowCommand() *cli.Command {
	return &cli.Command{
		Name:  "workflow",
		Usage: "Show workflow information",
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "show available workflows",
				Action: workflowCommandListAction,
			},
			{
				Name:   "show",
				Usage:  "show details of current workflow",
				Action: workflowCommandShowAction,
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name: "workflow-id",
					},
				},
			},
		},
	}
}

func workflowCommandListAction(_ context.Context, cmd *cli.Command) error {
	for _, workflowType := range workflow.SupportedTypes {
		_, _ = fmt.Fprintf(cmd.Writer, "%v\n", workflowType.Id())
	}

	return nil
}

func workflowPrintItem(writer io.Writer, name string, value any) {
	if value == nil {
		_, _ = fmt.Fprintf(writer, "%v: -\n", name)
	} else if tool, ok := value.(internal.Tool); ok {
		_, _ = fmt.Fprintf(writer, "%v: %v\n", name, tool.Id())
	} else if adaptor, ok := value.(workflow.Adaptor); ok {
		_, _ = fmt.Fprintf(writer, "%v: %v\n", name, adaptor.Details())
	} else {
		_, _ = fmt.Fprintf(writer, "%v: ?\n", name)
	}
}

func workflowPrintWorkflow(writer io.Writer, workflow internal.Workflow) {
	_, _ = fmt.Fprintf(writer, "Workflow: %v\n---\n", workflow.Name)

	workflowPrintItem(writer, "Configurator", workflow.Configurator)
	workflowPrintItem(writer, "Builder", workflow.Builder)
	workflowPrintItem(writer, "Formatter", workflow.Formatter)
	workflowPrintItem(writer, "Linter", workflow.Linter)
	workflowPrintItem(writer, "Runner", workflow.Runner)
	workflowPrintItem(writer, "Tester", workflow.Tester)
	workflowPrintItem(writer, "Dependency", workflow.DependencyManager)
}

func workflowCommandShowAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)

	if workflowId := cmd.StringArg("workflow-id"); workflowId != "" {
		for _, workflowType := range workflow.SupportedTypes {
			if workflowType.Id() == workflowId {
				config := workflow.Config{
					Directory:             cmdContext.Config.RootDirectory,
					IntermediateDirectory: cmdContext.Config.BuildDirectory,
				}

				workflowPrintWorkflow(
					cmd.Writer,
					workflowType.Create(config, cmdContext.Tools).Workflow,
				)
			}
		}
	} else {
		workflowPrintWorkflow(cmd.Writer, cmdContext.Project.Workflow)
	}

	return nil
}
