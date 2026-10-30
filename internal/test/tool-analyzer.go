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

// AnalyzerBuilder is a function that builds an analyzer for testing.
type AnalyzerBuilder func(executable func() (*internal.Executable, error)) internal.ProjectAnalyzer

// InvokeAnalyzeFiles run a test that verifies call AnalyzeFiles on an analyzer.
func InvokeAnalyzeFiles(
	t *testing.T,
	builder AnalyzerBuilder,
	options internal.ProjectAnalyzerOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectAnalyzer {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectAnalyzer) error {
			return tool.AnalyzeFiles(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeAnalyzeFilesSuccess run a test that verifies call AnalyzeFiles on an analyzer that should invoke the analyzer
// with expected arguments without error.
func InvokeAnalyzeFilesSuccess(
	t *testing.T,
	builder AnalyzerBuilder,
	options internal.ProjectAnalyzerOptions,
	expectedArgs []string,
) internal.ProjectAnalyzer {
	return InvokeAnalyzeFiles(
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

// InvokeAnalyzeFilesFailed run a test that verifies call AnalyzeFiles on an analyzer that should invoke the analyzer
// with expected arguments that return an error.
func InvokeAnalyzeFilesFailed(
	t *testing.T,
	builder AnalyzerBuilder,
	options internal.ProjectAnalyzerOptions,
	expectedArgs []string,
) internal.ProjectAnalyzer {
	expectedErr := errors.New("failed")

	return InvokeAnalyzeFiles(
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

// InvokeAnalyzeFilesError run a test that verifies call AnalyzeFiles on an analyzer that should not invoke analyzer
// and return an error. It should be used when some check before invoking the analyzer fails.
func InvokeAnalyzeFilesError(
	t *testing.T,
	builder AnalyzerBuilder,
	options internal.ProjectAnalyzerOptions,
	errorCheck func(err error),
) internal.ProjectAnalyzer {
	return InvokeAnalyzeFiles(
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

// InvokeAnalyzeFilesReporter run a test that verifies call AnalyzeFiles on an analyzer that should pass output to
// the reporter.
func InvokeAnalyzeFilesReporter(
	t *testing.T,
	builder AnalyzerBuilder,
	directory string,
	analyzeReporter internal.FileIssueReporter,
	expectedArgs []string,
	output string,
	errorCheck func(err error),
) internal.ProjectAnalyzer {
	defer func() { _ = analyzeReporter.Finish() }()

	return InvokeAnalyzeFiles(
		t,
		builder,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Output: internal.AnalyzeOutputOptions{
				Reporter: analyzeReporter,
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

// InvokeAnalyzeFilesOutputFile run a test that verifies call AnalyzeFiles on an analyzer that will output a file.
func InvokeAnalyzeFilesOutputFile(
	t *testing.T,
	builder AnalyzerBuilder,
	directory string,
	expectedArgs []string,
	filenames *[]string,
) internal.ProjectAnalyzer {
	const content = "ok"

	outputFilename := t.TempDir() + "/report.txt"

	file, err := os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	return InvokeAnalyzeFiles(
		t,
		builder,
		internal.ProjectAnalyzerOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory: directory,
			},
			Filenames: filenames,
			Output: internal.AnalyzeOutputOptions{
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
