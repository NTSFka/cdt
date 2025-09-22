package command

import (
	"cdt/internal/workflow"
	"context"
	"fmt"

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
		},
	}
}

func workflowCommandListAction(ctx context.Context, cmd *cli.Command) error {
	for _, workflowType := range workflow.SupportedTypes {
		_, _ = fmt.Fprintf(cmd.Writer, "%v\n", workflowType.Id())
	}

	return nil
}
