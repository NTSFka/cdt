package cli

import (
	"cdt/internal"
	. "cdt/pkg"
	"os"
	"testing"
)

func runMain(args ...string) error {
	return RunMain(RunContext{
		Args:   append([]string{"cdt"}, args...),
		Input:  os.Stdin,
		Output: os.Stdout,
		Error:  os.Stderr,
	})
}

// Check if tool exists in the current environment
func checkTool(t *testing.T, toolName string) {
	if executable := internal.FindExecutable(toolName); executable == nil {
		t.Skipf("unable to find tool: %v", toolName)
	}
}
