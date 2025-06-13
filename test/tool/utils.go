package tool

import (
	"cdt/internal"
	"testing"
)

// Check if tool exists in the current environment
func checkTool(t *testing.T, toolName string) {
	if executable := internal.FindExecutable(toolName); executable == nil {
		t.Skipf("unable to find tool: %v", toolName)
	}
}
