package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// RunnerBuilder is a function that builds a project runner for testing.
type RunnerBuilder func(executable func() (*internal.Executable, error)) internal.ProjectRunner

// InvokeRunTarget runs the RunTarget command and checks the result.
func InvokeRunTarget(
	t *testing.T,
	builder RunnerBuilder,
	options internal.ProjectRunnerOptions,
	target string,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectRunner {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectRunner) error {
			return tool.RunTarget(t.Context(), options, target)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeRunTargetSuccess run a test that verifies call RunTarget on a runner that should invoke the runner
// with expected arguments without error.
func InvokeRunTargetSuccess(
	t *testing.T,
	builder RunnerBuilder,
	options internal.ProjectRunnerOptions,
	target string,
	expectedArgs []string,
) internal.ProjectRunner {
	return InvokeRunTarget(
		t,
		builder,
		options,
		target,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

// InvokeRunTargetFailed run a test that verifies call RunTarget on a runner that should invoke the runner
// with expected arguments that return an error.
func InvokeRunTargetFailed(
	t *testing.T,
	builder RunnerBuilder,
	options internal.ProjectRunnerOptions,
	target string,
	expectedArgs []string,
) internal.ProjectRunner {
	expectedErr := errors.New("failed")

	return InvokeRunTarget(
		t,
		builder,
		options,
		target,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}
