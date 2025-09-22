package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runWorkflow(ctx context.Context, args ...string) error {
	return test.RunCommand(ctx, NewWorkflowCommand(), internal.Context{}, args...)
}

func TestWorkflowList(t *testing.T) {
	err := runWorkflow(context.Background(), "list")

	assert.NoError(t, err)
}
