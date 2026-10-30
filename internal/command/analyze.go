package command

import (
	"cdt/internal"
	"cdt/internal/reporter"
	"context"
	"errors"
	"fmt"
	"io"

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
		var err error

		analyzer, err = getTool[internal.ProjectAnalyzer](
			cmd.String("tool"),
			cmdContext.Tools,
			"analyzing",
		)
		if err != nil {
			return err
		}
	}

	if analyzer == nil {
		return errors.New("project doesn't support analyzing")
	}

	outputOptions := ParseOptionOutput(cmd.String("output"))

	output, err := buildOutput(outputOptions.Filename)
	if err != nil {
		return err
	}
	defer func() { _ = output.Close() }()

	rep, err := buildAnalyzeReporter(outputOptions.Format, output)
	if err != nil {
		return err
	}

	options := internal.ProjectAnalyzerOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
		Output: internal.AnalyzeOutputOptions{
			Output:   output,
			Reporter: rep,
		},
	}

	if files := cmd.StringArgs("files"); len(files) > 0 {
		options.Filenames = &files
	}

	if err := analyzer.AnalyzeFiles(ctx, options); err != nil {
		return fmt.Errorf("%w", err)
	}

	if rep != nil {
		return rep.Finish()
	}

	return nil
}

func buildAnalyzeReporter(format string, output io.Writer) (internal.FileIssueReporter, error) {
	switch format {
	case "":
		fallthrough
	case "raw":
		return nil, nil // nolint: nilnil
	case "simple":
		return reporter.NewSimpleSarifReporter(output), nil
	case "sarif":
		return reporter.NewFileSarifReporter(output), nil
	}

	return nil, fmt.Errorf("unsupported report format: %s", format)
}
