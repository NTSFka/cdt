package test

import (
	"cdt/internal"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BuilderBuilder is a function that builds a project builder for testing.
type BuilderBuilder func(executable func() (*internal.Executable, error)) internal.ProjectBuilder

// InvokeBuildTargets runs the build command and checks the result.
func InvokeBuildTargets(
	t *testing.T,
	builder BuilderBuilder,
	options internal.ProjectBuilderOptions,
	execMockSetup ExecMockSetup,
	resultCheck func(err error),
) internal.ProjectBuilder {
	return InvokeTool(
		t,
		builder,
		func(tool internal.ProjectBuilder) error {
			return tool.BuildTargets(t.Context(), options)
		},
		execMockSetup,
		resultCheck,
	)
}

// InvokeBuildTargetsSuccess runs a test that verifies call BuildTargets on the builder that should invoke
// the builder with expected arguments without error.
func InvokeBuildTargetsSuccess(
	t *testing.T,
	builder BuilderBuilder,
	options internal.ProjectBuilderOptions,
	expectedArgs []string,
) internal.ProjectBuilder {
	return InvokeBuildTargets(
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

// InvokeBuildTargetsFailed runs a test that verifies call BuildTargets on the builder that should invoke
// the builder with expected arguments with error.
func InvokeBuildTargetsFailed(
	t *testing.T,
	builder BuilderBuilder,
	options internal.ProjectBuilderOptions,
	expectedArgs []string,
) internal.ProjectBuilder {
	expectedErr := errors.New("failed")

	return InvokeBuildTargets(
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
