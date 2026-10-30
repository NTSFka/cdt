package command

import (
	"cdt/internal/reporter"
	"context"
	"errors"
	"fmt"
	"io"

	"cdt/internal"

	"github.com/urfave/cli/v3"
)

func NewLintCommand() *cli.Command {
	return &cli.Command{
		Name:   "lint",
		Usage:  "Lint the project",
		Action: lintCommandAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "tool",
				Usage: "Use specific linter tool",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "The lint report output format in form: <format>[:<filename>].",
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

func lintCommandAction(ctx context.Context, cmd *cli.Command) error {
	cmdContext := ctx.Value("context").(internal.Context)
	linter := cmdContext.Project.Workflow.Linter

	if cmd.IsSet("tool") {
		var err error

		linter, err = getTool[internal.ProjectLinter](
			cmd.String("tool"),
			cmdContext.Tools,
			"linting",
		)
		if err != nil {
			return err
		}
	}

	if linter == nil {
		return errors.New("project doesn't support linting")
	}

	outputOptions := ParseOptionOutput(cmd.String("output"))

	output, err := buildOutput(outputOptions.Filename)
	if err != nil {
		return err
	}
	defer func() { _ = output.Close() }()

	rep, err := buildLintReporter(outputOptions.Format, output)
	if err != nil {
		return err
	}

	options := internal.ProjectLinterOptions{
		ProjectInfo: cmdContext.Project.Info,
		ExtraArgs:   cmd.Args().Tail(),
		Output: internal.LintOutputOptions{
			Output:   output,
			Reporter: rep,
		},
	}

	if files := cmd.StringArgs("files"); len(files) > 0 {
		options.Filenames = &files
	}

	if err := linter.LintFiles(ctx, options); err != nil {
		return fmt.Errorf("%w", err)
	}

	if rep != nil {
		return rep.Finish()
	}

	return nil
}

func buildLintReporter(format string, output io.Writer) (internal.FileIssueReporter, error) {
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
