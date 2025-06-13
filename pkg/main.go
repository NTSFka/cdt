package pkg

import (
	. "cdt/internal"
	"cdt/internal/command"
	"cdt/internal/project"
	"cdt/internal/tool"
	"context"
	"github.com/urfave/cli/v3"
	"io"
	"os"
)

// InitTools initializes all supported tools on the system
func initTools() Tools {
	return Tools{
		tool.DetectClangFormat(nil),
		tool.DetectClangTidy(nil),
		tool.DetectCMake(),
		tool.DetectCTest(),
	}
}

func detectProject(rootDirectory string, buildDirectory string, tools Tools) Project {
	// CMake
	if p := project.DetectCMakeProject(rootDirectory, buildDirectory, tools); p != nil {
		return *p
	}

	return MakeProject(rootDirectory, buildDirectory, &EmptyProjectStructureProvider{}, Workflow{})
}

type RunContext struct {
	Args   []string
	Input  io.Reader
	Output io.Writer
	Error  io.Writer
}

// NewRunContext create standard run context
func NewRunContext() RunContext {
	return RunContext{
		Args:   os.Args,
		Input:  os.Stdin,
		Output: os.Stdout,
		Error:  os.Stderr,
	}
}

func RunMain(runCtx RunContext) error {
	tools := initTools()

	cmd := &cli.Command{
		Name:                  "cdt",
		Usage:                 "A common developer tool",
		Version:               "0.1.0",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"d"},
				Usage:   "project directory",
				Value:   ".",
			},
			&cli.StringFlag{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "build directory",
				Value:   "build",
			},
		},
		Commands: []*cli.Command{
			&command.ProjectCommand,
			&command.ToolCommand,
			&command.ConfigureCommand,
			&command.BuildCommand,
			&command.FormatCommand,
			&command.TestCommand,
			&command.LintCommand,
			&command.RunCommand,
		},
		Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
			projectPath := command.String("directory")
			buildDirectory := command.String("build")

			proj := detectProject(projectPath, buildDirectory, tools)

			c := Context{
				Project: proj,
				Tools:   tools,
			}

			return context.WithValue(ctx, "context", c), nil
		},
	}

	return cmd.Run(context.Background(), runCtx.Args)
}
