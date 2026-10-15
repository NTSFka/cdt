package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// ExecMockSetup is a function that sets up an executable mock.
type ExecMockSetup func(exec *Executable, executablePath string)

// InvokeTool runs a test that verifies a tool invocation.
//
// The function tests a tool by preparing executable mock, invoking the tool, and checking the result.
func InvokeTool[T any](
	t *testing.T,
	buildTool func(executable func() (*internal.Executable, error)) T,
	invokeTool func(tool T) error,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) T {
	const executablePath = "/bin/tool"

	exec := NewExecutable(t)

	tool := buildTool(exec.LazyExecutable(executablePath))

	execMockSetup(exec, executablePath)

	err := invokeTool(tool)
	resultCheck(err)

	exec.AssertExpectations(t)

	return tool
}

// InvokeToolResult runs a test that verifies a tool invocation that returns some result.
//
// The function tests a tool by preparing executable mock, invoking the tool, and checking the result.
func InvokeToolResult[T any, U any](
	t *testing.T,
	buildTool func(executable func() (*internal.Executable, error)) T,
	invokeTool func(tool T) (U, error),
	execMockSetup ExecMockSetup,
	resultCheck func(result U, err error),
) T {
	const executablePath = "/bin/tool"

	exec := NewExecutable(t)

	tool := buildTool(exec.LazyExecutable(executablePath))

	execMockSetup(exec, executablePath)

	result, err := invokeTool(tool)
	resultCheck(result, err)

	exec.AssertExpectations(t)

	return tool
}

// ToolBuilder is a function that builds a tool for testing.
type ToolBuilder func(executable func() (*internal.Executable, error)) internal.Tool

// InvokeRun runs a test that verifies Run call on a tool.
func InvokeRun(
	t *testing.T,
	builder ToolBuilder,
	runOptions internal.RunOptions,
	args []string,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.Tool {
	return InvokeTool(
		t,
		builder,
		func(tool internal.Tool) error {
			return tool.Run(t.Context(), runOptions, args)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeRunSuccess run a test that verifies Run call on a tool without error.
func InvokeRunSuccess(
	t *testing.T,
	builder ToolBuilder,
	runOptions internal.RunOptions,
	args []string,
	expectedArgs []string,
) internal.Tool {
	return InvokeRun(
		t,
		builder,
		runOptions,
		args,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

// InvokeRunFailed runs a test that verifies Run call on a tool that will fail.
func InvokeRunFailed(
	t *testing.T,
	builder ToolBuilder,
	runOptions internal.RunOptions,
	args []string,
	expectedArgs []string,
) internal.Tool {
	expectedErr := errors.New("failed")

	return InvokeRun(
		t,
		builder,
		runOptions,
		args,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

// InvokeRunForProject runs a test that verifies RunForProject call on a tool.
func InvokeRunForProject(
	t *testing.T,
	builder ToolBuilder,
	info internal.ProjectInfo,
	args []string,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.Tool {
	return InvokeTool(
		t,
		builder,
		func(tool internal.Tool) error {
			return tool.RunForProject(t.Context(), info, args)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeRunForProjectSuccess run a test that verifies RunForProject call on a tool without error.
func InvokeRunForProjectSuccess(
	t *testing.T,
	builder ToolBuilder,
	info internal.ProjectInfo,
	args []string,
	expectedArgs []string,
) internal.Tool {
	return InvokeRunForProject(
		t,
		builder,
		info,
		args,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

// InvokeRunForProjectFailed runs a test that verifies RunForProject call on a tool that will fail.
func InvokeRunForProjectFailed(
	t *testing.T,
	builder ToolBuilder,
	info internal.ProjectInfo,
	args []string,
	expectedArgs []string,
) internal.Tool {
	expectedErr := errors.New("failed")

	return InvokeRunForProject(
		t,
		builder,
		info,
		args,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}
