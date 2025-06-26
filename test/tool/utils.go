package tool

import (
	"cdt/internal"
	"testing"
)

// Check if a tool exists in the given environment
func checkTool(t *testing.T, environment internal.Environment, toolName string) {
	if executable := environment.FindExecutable(toolName); executable == nil {
		t.Skipf("unable to find tool: %v", toolName)
	}
}
