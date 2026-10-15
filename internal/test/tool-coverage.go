package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CoverageCollectorBuilder is a function that builds a project runner for testing.
type CoverageCollectorBuilder func(executable func() (*internal.Executable, error)) internal.ProjectCoverageCollector

// InvokeCollectCoverage runs the build command and checks the result.
func InvokeCollectCoverage(
	t *testing.T,
	builder CoverageCollectorBuilder,
	options internal.ProjectCoverageCollectorOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectCoverageCollector {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectCoverageCollector) error {
			return tool.CollectCoverage(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeCollectCoverageSuccess runs a test that verifies call CollectCoverage on the collector that should invoke
// the collector with expected arguments without error.
func InvokeCollectCoverageSuccess(
	t *testing.T,
	builder CoverageCollectorBuilder,
	options internal.ProjectCoverageCollectorOptions,
	expectedArgs []string,
) internal.ProjectCoverageCollector {
	return InvokeCollectCoverage(
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

// InvokeCollectCoverageFailed runs a test that verifies call CollectCoverage on the collector that should invoke
// the collector with expected arguments with error.
func InvokeCollectCoverageFailed(
	t *testing.T,
	builder CoverageCollectorBuilder,
	options internal.ProjectCoverageCollectorOptions,
	expectedArgs []string,
) internal.ProjectCoverageCollector {
	expectedErr := errors.New("failed")

	return InvokeCollectCoverage(
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

// InvokeCollectCoverageError run a test that verifies call CollectCoverage on a collector that should not invoke
// collector and return an error. It should be used when some check before invoking the collector fails.
func InvokeCollectCoverageError(
	t *testing.T,
	builder CoverageCollectorBuilder,
	options internal.ProjectCoverageCollectorOptions,
	errorCheck func(err error),
) internal.ProjectCoverageCollector {
	return InvokeCollectCoverage(
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

// InvokeCollectCoverageOutputFormatUnsupported run a test that verifies call CollectCoverage on a collector for
// unsupported output format.
func InvokeCollectCoverageOutputFormatUnsupported(
	t *testing.T,
	builder CoverageCollectorBuilder,
	format internal.CoverageReportFormat,
	pattern *string,
) internal.ProjectCoverageCollector {
	return InvokeCollectCoverageError(
		t,
		builder,
		internal.ProjectCoverageCollectorOptions{
			Output: internal.OutputOptions[internal.CoverageReportFormat]{
				Format: format,
			},
			Pattern: pattern,
		},
		func(err error) {
			require.EqualError(t, err, "unsupported coverage report format: "+string(format))
		},
	)
}
