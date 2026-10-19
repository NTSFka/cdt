package test

import (
	"cdt/internal"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LinterBuilder is a function that builds a linter for testing.
type LinterBuilder func(executable func() (*internal.Executable, error)) internal.ProjectLinter

// InvokeLintFiles run a test that verifies call LintFiles on a linter.
func InvokeLintFiles(
	t *testing.T,
	builder LinterBuilder,
	options internal.ProjectLinterOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectLinter {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectLinter) error {
			return tool.LintFiles(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeLintFilesSuccess run a test that verifies call LintFiles on a linter that should invoke the linter with
// expected arguments without error.
func InvokeLintFilesSuccess(
	t *testing.T,
	builder LinterBuilder,
	options internal.ProjectLinterOptions,
	expectedArgs []string,
) internal.ProjectLinter {
	return InvokeLintFiles(
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

// InvokeLintFilesFailed run a test that verifies call LintFiles on a linter that should invoke the linter with
// expected arguments that return an error.
func InvokeLintFilesFailed(
	t *testing.T,
	builder LinterBuilder,
	options internal.ProjectLinterOptions,
	expectedArgs []string,
) internal.ProjectLinter {
	expectedErr := errors.New("failed")

	return InvokeLintFiles(
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

// InvokeLintFilesError run a test that verifies call LintFiles on a linter that should not invoke linter and
// return an error. It should be used when some check before invoking the linter fails.
func InvokeLintFilesError(
	t *testing.T,
	builder LinterBuilder,
	options internal.ProjectLinterOptions,
	errorCheck func(err error),
) internal.ProjectLinter {
	return InvokeLintFiles(
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

// InvokeLintFilesOutputFormatCheck run a test that verifies call LintFiles on a linter that should invoke the linter
// with expected arguments for a given output format.
func InvokeLintFilesOutputFormatCheck(
	t *testing.T,
	builder LinterBuilder,
	directory string,
	format internal.LintReportFormat,
	expectedArgs []string,
	filenames *[]string,
) internal.ProjectLinter {
	return InvokeLintFilesSuccess(
		t,
		builder,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Filenames: filenames,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: format,
			},
		},
		expectedArgs,
	)
}

// InvokeLintFilesOutputFormatUnsupported run a test that verifies call LintFiles on a linter for unsupported output
// format.
func InvokeLintFilesOutputFormatUnsupported(
	t *testing.T,
	builder LinterBuilder,
	directory string,
	format internal.LintReportFormat,
	filenames *[]string,
) internal.ProjectLinter {
	return InvokeLintFilesError(
		t,
		builder,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Filenames: filenames,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: format,
			},
		},
		func(err error) {
			require.EqualError(t, err, "unsupported report format: "+string(format))
		},
	)
}

// InvokeLintFilesOutputFile run a test that verifies call LintFiles on a linter that will output a file.
func InvokeLintFilesOutputFile(
	t *testing.T,
	builder LinterBuilder,
	directory string,
	expectedArgs []string,
	filenames *[]string,
) internal.ProjectLinter {
	const content = "ok"

	outputFilename := t.TempDir() + "/report.txt"

	return InvokeLintFiles(
		t,
		builder,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Filenames: filenames,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Filename: &outputFilename,
			},
		},
		func(exec *Executable, executablePath string) {
			exec.OnRunOutput(executablePath, expectedArgs, content).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
			require.FileExists(t, outputFilename)

			outputContent, err := os.ReadFile(outputFilename)
			require.NoError(t, err)
			assert.Equal(t, content, string(outputContent))
		},
	)
}
