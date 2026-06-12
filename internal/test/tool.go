package test

import (
	"cdt/internal"
	"cdt/internal/tool"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ToolDetector is a function that tries to detect a tool.
type ToolDetector func(context.Context, tool.DetectOptions) internal.Tool

// RunDetectTool run a test that verifies function for detecting a tool.
//
// The detector function is used to detect a tool = tested function.
// The envMockSetup function is used to set up the environment mock.
// The resultCheck function is used to check the result of the detection.
// The toolPaths are user mapping of names to paths.
func RunDetectTool(
	t *testing.T,
	detector ToolDetector,
	envMockSetup func(env *Environment),
	resultCheck func(tool internal.Tool),
	toolPaths map[string]string,
) internal.Tool {
	env := NewEnvironment(t)

	envMockSetup(env)

	detectedTool := detector(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  toolPaths,
	})
	resultCheck(detectedTool)

	env.AssertExpectations(t)

	return detectedTool
}

// RunDetectToolLastFound run a test that verifies function for detecting a tool by checking testPaths.
// It should succeed for the last path in testPaths.
//
// The detector function is used to detect a tool = tested function.
// The testPaths are the paths that the function should check. The last path is the one that should be found.
func RunDetectToolLastFound(
	t *testing.T,
	detector ToolDetector,
	testPaths []string,
) internal.Tool {
	return RunDetectTool(
		t,
		detector,
		func(env *Environment) {
			for idx, path := range testPaths {
				exe := env.OnFindExecutable(path)

				// Last path is the one that should be found.
				if idx == len(testPaths)-1 {
					exe.Return(env.NewExecutable("/bin/tool"), nil)
				} else {
					exe.Return(nil, nil)
				}
			}
		},
		func(tool internal.Tool) {
			require.NotNil(t, tool)
			assert.True(t, tool.IsAvailable())

			executable := tool.Executable()
			require.NotNil(t, executable)
			assert.Equal(t, "/bin/tool", executable.Path)
		},
		nil,
	)
}

// RunDetectToolNotFound run a test that verifies function for detecting a tool by checking testPaths.
// It should fail for all testPaths.
//
// The detector function is used to detect a tool = tested function.
// The testPaths are the paths that the function should check.
func RunDetectToolNotFound(
	t *testing.T,
	detector ToolDetector,
	testPaths []string,
) internal.Tool {
	return RunDetectTool(
		t,
		detector,
		func(env *Environment) {
			for _, path := range testPaths {
				env.OnFindExecutable(path).Return(nil, nil)
			}
		},
		func(tool internal.Tool) {
			require.NotNil(t, tool)
			assert.False(t, tool.IsAvailable())
			assert.Nil(t, tool.Executable())
		},
		nil,
	)
}

// RunDetectToolConfig run a test that verifies function for detecting a tool with custom tool mapping.
//
// The detector function is used to detect a tool = tested function.
// The searchName is the name of the tool in the tool mapping.
func RunDetectToolConfig(t *testing.T, detector ToolDetector, searchName string) internal.Tool {
	const toolPath = "tool-path"

	return RunDetectTool(
		t,
		detector,
		func(env *Environment) {
			env.OnFindExecutable(toolPath).
				Return(env.NewExecutable("/bin/tool"), nil)
		},
		func(tool internal.Tool) {
			assert.True(t, tool.IsAvailable())
			executable := tool.Executable()

			require.NotNil(t, executable)
			assert.Equal(t, "/bin/tool", executable.Path)
		},
		map[string]string{searchName: toolPath},
	)
}

// LinterBuilder is a function that builds a linter for testing.
type LinterBuilder func(executable func() (*internal.Executable, error)) internal.ProjectLinter

// RunLintFiles run a test that verifies call LintFiles on a linter.
//
// The builder function is used to build tested linter.
// The options are the options to pass to the LintFiles call.
// The execMockSetup function is used to set up the executable mock.
// The resultCheck function is used to check the result of the LintFiles call.
func RunLintFiles(
	t *testing.T,
	builder LinterBuilder,
	options internal.ProjectLinterOptions,
	executablePath string,
	execMockSetup func(exec *Executable),
	resultCheck func(err error),
) internal.ProjectLinter {
	exec := NewExecutable(t)

	linter := builder(exec.LazyExecutable(executablePath))

	execMockSetup(exec)

	err := linter.LintFiles(t.Context(), options)
	resultCheck(err)

	exec.AssertExpectations(t)

	return linter
}

// RunLintFilesInvoked run a test that verifies call LintFiles on a linter that should invoke the linter with expected
// arguments.
//
// The builder function is used to build tested linter.
// The options are the options to pass to the LintFiles call.
// The expectedArgs are the expected arguments the linter should be invoked with.
// The resultCheck function is used to check the result of the LintFiles call.
func RunLintFilesInvoked(
	t *testing.T,
	builder LinterBuilder,
	options internal.ProjectLinterOptions,
	expectedArgs []string,
	resultCheck func(err error),
) internal.ProjectLinter {
	const executablePath = "lint"

	return RunLintFiles(
		t,
		builder,
		options,
		executablePath,
		func(exec *Executable) {
			// Check that the linter was invoked with the expected arguments.
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		resultCheck,
	)
}

// RunLintFilesSuccess run a test that verifies call LintFiles on a linter that should invoke the linter with expected
// arguments without error.
//
// The builder function is used to build tested linter.
// The options are the options to pass to the LintFiles call.
// The expectedArgs are the expected arguments the linter should be invoked with.
func RunLintFilesSuccess(
	t *testing.T,
	builder LinterBuilder,
	options internal.ProjectLinterOptions,
	expectedArgs []string,
) internal.ProjectLinter {
	return RunLintFilesInvoked(t, builder, options, expectedArgs, func(err error) {
		require.NoError(t, err)
	})
}

// RunLintFilesOutputFormatCheck run a test that verifies call LintFiles on a linter that should invoke the linter
// with expected arguments for a given output format.
//
// The builder function is used to build tested linter.
// The format is the output format.
// The expectedArgs are the expected arguments the linter should be invoked with.
// The filenames are the filenames to lint. Nil value if for "all files".
func RunLintFilesOutputFormatCheck(
	t *testing.T,
	builder LinterBuilder,
	format internal.LintReportFormat,
	expectedArgs []string,
	filenames *[]string,
) internal.ProjectLinter {
	return RunLintFilesSuccess(
		t,
		builder,
		internal.ProjectLinterOptions{
			Filenames: filenames,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: format,
			},
		},
		expectedArgs,
	)
}

// RunLintFilesOutputFormatUnsupported run a test that verifies call LintFiles on a linter for unsupported output
// format.
//
// The builder function is used to build tested linter.
// The format is the output format (unsupported).
// The filenames are the filenames to lint. Nil value if for "all files".
func RunLintFilesOutputFormatUnsupported(
	t *testing.T,
	builder LinterBuilder,
	format internal.LintReportFormat,
	filenames *[]string,
) internal.ProjectLinter {
	return RunLintFiles(
		t,
		builder,
		internal.ProjectLinterOptions{
			Filenames: filenames,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: format,
			},
		},
		"linter",
		func(exec *Executable) {
			// No expectations needed.
		},
		func(err error) {
			require.EqualError(t, err, "unsupported report format: "+string(format))
		},
	)
}

// RunLintFilesOutputFile run a test that verifies call LintFiles on a linter that will output a file.
//
// The builder function is used to build tested linter.
// The expectedArgs are the expected arguments the linter should be invoked with.
// The filenames are the filenames to lint. Nil value if for "all files".
func RunLintFilesOutputFile(
	t *testing.T,
	builder LinterBuilder,
	expectedArgs []string,
	filenames *[]string,
) internal.ProjectLinter {
	const (
		executablePath = "linter"
		content        = "ok"
	)

	outputFilename := t.TempDir() + "/report.txt"

	return RunLintFiles(
		t,
		builder,
		internal.ProjectLinterOptions{
			Filenames: filenames,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Filename: &outputFilename,
			},
		},
		executablePath,
		func(exec *Executable) {
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
