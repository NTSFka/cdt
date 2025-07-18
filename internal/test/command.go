package test

import (
	"cdt/internal"
	"context"
	"github.com/urfave/cli/v3"
)

func RunCommand(command *cli.Command, c internal.Context, args ...string) error {
	ctx := context.WithValue(context.Background(), "context", c) //nolint:staticcheck

	return command.Run(ctx, append([]string{"app"}, args...))
}
