package command

import (
	"context"
	"errors"
	"fmt"

	"cdt/internal"

	"github.com/urfave/cli/v3"
)

func NewAnalyzeCommand() *cli.Command {
	return &cli.Command{
		Name:   "analyze",
		Usage:  "Analyze the project",
		Action: analyzeCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific analyzer tool",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "The analyze report output format in form: <format>[:<filename>].",
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name: "files",
				Min:  0,
				Max:  -1,
			},
		},
	}
}

func analyzeCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	analyzer := cmdContext.Project.Workflow.Analyzer

	if cmd.IsSet("tool") {
		toolName := cmd.String("tool")
		tool := cmdContext.Tools.Get(toolName)

		if tool == nil {
			return fmt.Errorf("tool '%s' not found", toolName)
		}

		linterTool, ok := tool.(internal.ProjectAnalyzer)

		if ok {
			analyzer = linterTool
		} else {
			return fmt.Errorf("tool '%s' doesn't support analyzing", toolName)
		}
	}

	if analyzer == nil {
		return errors.New("project doesn't support analyzing")
	}

	outputOptions := ParseOptionOutput[internal.AnalyzeReportFormat](
		cmd.String("output"),
		internal.AnalyzeReportFormatDefault,
	)

	options := internal.ProjectAnalyzerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
		Output:      outputOptions,
	}

	if files := cmd.StringArgs("files"); len(files) > 0 {
		options.Filenames = &files
	}

	if err := analyzer.AnalyzeFiles(ctx, options); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
