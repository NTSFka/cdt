package test

import (
	"cdt/internal"
	"errors"
	"io"
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

// InvokeLintFilesReporter run a test that verifies call LintFiles on a linter that should pass output to
// the reporter.
func InvokeLintFilesReporter(
	t *testing.T,
	builder LinterBuilder,
	directory string,
	linterReporter internal.FileIssueReporter,
	expectedArgs []string,
	output string,
	errorCheck func(err error),
) internal.ProjectLinter {
	defer func() { _ = linterReporter.Finish() }()

	return InvokeLintFiles(
		t,
		builder,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.LintOutputOptions{
				Reporter: linterReporter,
				Output:   io.Discard,
			},
		},
		func(exec *Executable, executablePath string) {
			exec.OnRunOutput(executablePath, expectedArgs, output).
				Return(nil)
		},
		errorCheck,
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

	file, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	return InvokeLintFiles(
		t,
		builder,
		internal.ProjectLinterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Filenames: filenames,
			Output: internal.LintOutputOptions{
				Output: file,
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
