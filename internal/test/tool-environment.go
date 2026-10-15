package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// EnvironmentProviderBuilder is a function that builds an environment provider for testing.
type EnvironmentProviderBuilder func(executable func() (*internal.Executable, error)) internal.EnvironmentProvider

// InvokeCreateEnvironment run a test that verifies call CreateEnvironment on an environment provider.
func InvokeCreateEnvironment(
	t *testing.T,
	builder EnvironmentProviderBuilder,
	directory string,
	parameter string,
	execMockSetup ExecMockSetup,
	resultCheck func(env internal.Environment, err error),
) internal.EnvironmentProvider {
	return InvokeToolResult(
		t,
		builder,
		func(tool internal.EnvironmentProvider) (internal.Environment, error) {
			return tool.CreateEnvironment(directory, parameter)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeCreateEnvironmentSuccess run a test that verifies call CreateEnvironment on an environment provider that
// should invoke the environment provider with expected arguments without error.
func InvokeCreateEnvironmentSuccess(
	t *testing.T,
	builder EnvironmentProviderBuilder,
	directory string,
	parameter string,
	expectedArgs []string,
) internal.EnvironmentProvider {
	return InvokeCreateEnvironment(
		t,
		builder,
		directory,
		parameter,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(nil)
		},
		func(env internal.Environment, err error) {
			require.NoError(t, err)
			require.NotNil(t, env)
		},
	)
}

// InvokeCreateEnvironmentFailed run a test that verifies call CreateEnvironment on an environment provider that
// should invoke the environment provider with expected arguments that return an error.
func InvokeCreateEnvironmentFailed(
	t *testing.T,
	builder EnvironmentProviderBuilder,
	directory string,
	parameter string,
	expectedArgs []string,
) internal.EnvironmentProvider {
	expectedErr := errors.New("failed")

	return InvokeCreateEnvironment(
		t,
		builder,
		directory,
		parameter,
		func(exec *Executable, executablePath string) {
			exec.OnRun(executablePath, expectedArgs).Return(expectedErr)
		},
		func(env internal.Environment, err error) {
			require.ErrorIs(t, err, expectedErr)
			require.Nil(t, env)
		},
	)
}
