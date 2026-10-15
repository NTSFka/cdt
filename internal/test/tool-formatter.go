package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// FormatterBuilder is a function that builds a formatter for testing.
type FormatterBuilder func(executable func() (*internal.Executable, error)) internal.ProjectFormatter

// InvokeFormatFiles run a test that verifies call FormatFiles on a formatter.
func InvokeFormatFiles(
	t *testing.T,
	builder FormatterBuilder,
	options internal.ProjectFormatterOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectFormatter {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectFormatter) error {
			return tool.FormatFiles(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeFormatFilesSuccess run a test that verifies call FormatFiles on a formatter that should invoke the formatter
// with expected arguments without error.
func InvokeFormatFilesSuccess(
	t *testing.T,
	builder FormatterBuilder,
	options internal.ProjectFormatterOptions,
	expectedArgs []string,
) internal.ProjectFormatter {
	return InvokeFormatFiles(
		t,
		builder,
		options,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

// InvokeFormatFilesFailed run a test that verifies call FormatFiles on a formatter that should invoke the formatter
// with expected arguments that return an error.
func InvokeFormatFilesFailed(
	t *testing.T,
	builder FormatterBuilder,
	options internal.ProjectFormatterOptions,
	expectedArgs []string,
) internal.ProjectFormatter {
	expectedErr := errors.New("failed")

	return InvokeFormatFiles(
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
