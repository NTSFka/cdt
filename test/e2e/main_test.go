package e2e_test

import (
	"cdt/internal"
	"context"
	"io"
	"os"
	"os/exec"
	"testing"
)

const BinaryName = "cdt"

// Check if a tool exists in the given environment.
func checkTool(
	t *testing.T,
	ctx context.Context,
	environment internal.Environment,
	toolName string,
) {
	if executable, _ := environment.FindExecutable(ctx, toolName); executable == nil {
		t.Skipf("unable to find tool: %v", toolName)
	}
}

func runCdtOutput(
	ctx context.Context,
	output io.Writer,
	outputError io.Writer,
	args ...string,
) error {
	cmd := exec.CommandContext(ctx, "./"+BinaryName, args...)
	cmd.Stdout = output
	cmd.Stderr = outputError

	return cmd.Run()
}

func runCdt(ctx context.Context, args ...string) error {
	return runCdtOutput(ctx, os.Stdout, os.Stderr, args...)
}

func TestMain(m *testing.M) {
	// Build the binary before running tests
	cmd := exec.Command("go", "build", BinaryName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	_ = os.Remove(BinaryName)
	os.Exit(code)
}
