package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ConfiguratorBuilder is a function that builds a project configurator for testing.
type ConfiguratorBuilder func(executable func() (*internal.Executable, error)) internal.ProjectConfigurator

// InvokeConfigure runs the configure command and checks the result.
func InvokeConfigure(
	t *testing.T,
	builder ConfiguratorBuilder,
	options internal.ProjectConfiguratorOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectConfigurator {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectConfigurator) error {
			return tool.Configure(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeConfigureSuccess runs a test that verifies call Configure on the configurator that should invoke
// the configurator with expected arguments without error.
func InvokeConfigureSuccess(
	t *testing.T,
	builder ConfiguratorBuilder,
	options internal.ProjectConfiguratorOptions,
	expectedArgs []string,
) internal.ProjectConfigurator {
	return InvokeConfigure(
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

// InvokeConfigureFailed runs a test that verifies call Configure on the configurator that should invoke
// the configurator with expected arguments with error.
func InvokeConfigureFailed(
	t *testing.T,
	builder ConfiguratorBuilder,
	options internal.ProjectConfiguratorOptions,
	expectedArgs []string,
) internal.ProjectConfigurator {
	expectedErr := errors.New("failed")

	return InvokeConfigure(
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
