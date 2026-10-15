package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// DependencyManagerBuilder is a function that builds a project dependency manager for testing.
type DependencyManagerBuilder func(executable func() (*internal.Executable, error)) internal.ProjectDependencyManager

// InvokeAddDependencies runs the build command and checks the result.
func InvokeAddDependencies(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectDependencyManager {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectDependencyManager) error {
			return tool.AddDependencies(t.Context(), options, dependencies, dev)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeAddDependenciesSuccess runs a test that verifies call AddDependencies on the dependency manager that should
// invoke the dependency manager with expected arguments without error.
func InvokeAddDependenciesSuccess(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	return InvokeAddDependencies(
		t,
		builder,
		options,
		dependencies,
		dev,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			assert.NoError(t, err)
		},
	)
}

// InvokeAddDependenciesFailed runs a test that verifies call AddDependencies on the dependency manager that should
// invoke the dependency manager with expected arguments with error.
func InvokeAddDependenciesFailed(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	expectedErr := errors.New("failed")

	return InvokeAddDependencies(
		t,
		builder,
		options,
		dependencies,
		dev,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

// InvokeRemoveDependencies runs the build command and checks the result.
func InvokeRemoveDependencies(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectDependencyManager {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectDependencyManager) error {
			return tool.RemoveDependencies(t.Context(), options, dependencies, dev)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeRemoveDependenciesSuccess runs a test that verifies call RemoveDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments without error.
func InvokeRemoveDependenciesSuccess(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	return InvokeRemoveDependencies(
		t,
		builder,
		options,
		dependencies,
		dev,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			assert.NoError(t, err)
		},
	)
}

// InvokeRemoveDependenciesFailed runs a test that verifies call RemoveDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments with error.
func InvokeRemoveDependenciesFailed(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	dev bool,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	expectedErr := errors.New("failed")

	return InvokeRemoveDependencies(
		t,
		builder,
		options,
		dependencies,
		dev,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

// InvokeUpdateDependencies runs the build command and checks the result.
func InvokeUpdateDependencies(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectDependencyManager {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectDependencyManager) error {
			return tool.UpdateDependencies(t.Context(), options, dependencies)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeUpdateDependenciesSuccess runs a test that verifies call UpdateDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments without error.
func InvokeUpdateDependenciesSuccess(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	return InvokeUpdateDependencies(
		t,
		builder,
		options,
		dependencies,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			assert.NoError(t, err)
		},
	)
}

// InvokeUpdateDependenciesFailed runs a test that verifies call UpdateDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments with error.
func InvokeUpdateDependenciesFailed(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	dependencies []string,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	expectedErr := errors.New("failed")

	return InvokeUpdateDependencies(
		t,
		builder,
		options,
		dependencies,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

// InvokeFetchDependencies runs the build command and checks the result.
func InvokeFetchDependencies(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	noDev bool,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectDependencyManager {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectDependencyManager) error {
			return tool.FetchDependencies(t.Context(), options, noDev)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeFetchDependenciesSuccess runs a test that verifies call FetchDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments without error.
func InvokeFetchDependenciesSuccess(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	noDev bool,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	return InvokeFetchDependencies(
		t,
		builder,
		options,
		noDev,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(err error) {
			assert.NoError(t, err)
		},
	)
}

// InvokeFetchDependenciesFailed runs a test that verifies call FetchDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments with error.
func InvokeFetchDependenciesFailed(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	noDev bool,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	expectedErr := errors.New("failed")

	return InvokeFetchDependencies(
		t,
		builder,
		options,
		noDev,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(err error) {
			require.ErrorIs(t, err, expectedErr)
		},
	)
}

// InvokeListDependencies runs the build command and checks the result.
func InvokeListDependencies(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectDependencyManager {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectDependencyManager) error {
			return tool.ListDependencies(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeListDependenciesSuccess runs a test that verifies call ListDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments without error.
func InvokeListDependenciesSuccess(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	return InvokeListDependencies(
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

// InvokeListDependenciesFailed runs a test that verifies call ListDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments with error.
func InvokeListDependenciesFailed(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	expectedErr := errors.New("failed")

	return InvokeListDependencies(
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

// InvokeAuditDependencies runs the build command and checks the result.
func InvokeAuditDependencies(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectDependencyManager {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectDependencyManager) error {
			return tool.AuditDependencies(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeAuditDependenciesSuccess runs a test that verifies call AuditDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments without error.
func InvokeAuditDependenciesSuccess(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	return InvokeAuditDependencies(
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

// InvokeAuditDependenciesFailed runs a test that verifies call AuditDependencies on the dependency manager that
// should invoke the dependency manager with expected arguments with error.
func InvokeAuditDependenciesFailed(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	expectedArgs []string,
) internal.ProjectDependencyManager {
	expectedErr := errors.New("failed")

	return InvokeAuditDependencies(
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

// InvokeAuditDependenciesError run a test that verifies call AuditDependencies on a dependency manager that should
// not invoke dependency manager and return an error. It should be used when some check before invoking the dependency
// manager fails.
func InvokeAuditDependenciesError(
	t *testing.T,
	builder DependencyManagerBuilder,
	options internal.ProjectDependencyManagerOptions,
	errorCheck func(err error),
) internal.ProjectDependencyManager {
	return InvokeAuditDependencies(
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
