package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TesterBuilder is a function that builds a project tester for testing.
type TesterBuilder func(executable func() (*internal.Executable, error)) internal.ProjectTester

// InvokeRunTests runs the build command and checks the result.
func InvokeRunTests(
	t *testing.T,
	builder TesterBuilder,
	options internal.ProjectTesterOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectTester {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectTester) error {
			return tool.RunTests(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeRunTestsSuccess runs a test that verifies call RunTests on the tester that should invoke
// the tester with expected arguments without error.
func InvokeRunTestsSuccess(
	t *testing.T,
	builder TesterBuilder,
	options internal.ProjectTesterOptions,
	expectedArgs []string,
) internal.ProjectTester {
	return InvokeRunTests(
		t,
		builder,
		options,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			assert.NoError(t, err)
		},
	)
}

// InvokeRunTestsFailed runs a test that verifies call RunTests on the tester that should invoke
// the tester with expected arguments with error.
func InvokeRunTestsFailed(
	t *testing.T,
	builder TesterBuilder,
	options internal.ProjectTesterOptions,
	expectedArgs []string,
) internal.ProjectTester {
	expectedErr := errors.New("failed")

	return InvokeRunTests(
		t,
		builder,
		options,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

// InvokeRunTestsError run a test that verifies call RunTests on a tester that should not invoke tester and
// return an error. It should be used when some check before invoking the tester fails.
func InvokeRunTestsError(
	t *testing.T,
	builder TesterBuilder,
	options internal.ProjectTesterOptions,
	errorCheck func(err error),
) internal.ProjectTester {
	return InvokeRunTests(
		t,
		builder,
		options,
		func(exec *Executable, executablePath string) {
			// No invocation
		},
		func(err error) {
			require.Error(t, err)
			errorCheck(err)
		},
	)
}

// InvokeRunTestsOutputFormatUnsupported run a test that verifies call RunTests on a tester for unsupported output
// format.
func InvokeRunTestsOutputFormatUnsupported(
	t *testing.T,
	builder TesterBuilder,
	format internal.TestsReportFormat,
	pattern *string,
) internal.ProjectTester {
	return InvokeRunTestsError(
		t,
		builder,
		internal.ProjectTesterOptions{
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format: format,
			},
			Pattern: pattern,
		},
		func(err error) {
			require.EqualError(t, err, "unsupported report format: "+string(format))
		},
	)
}
