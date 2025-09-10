package tool

import (
	"cdt/internal"
	"context"
	"testing"
)

// Check if a tool exists in the given environment
func checkTool(t *testing.T, ctx context.Context, environment internal.Environment, toolName string) {
	if executable := environment.FindExecutable(ctx, toolName); executable == nil {
		t.Skipf("unable to find tool: %v", toolName)
	}
}
