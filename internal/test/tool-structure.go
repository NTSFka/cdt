package test

import (
	"cdt/internal"
	"testing"
)

// StructureProviderBuilder is a function that builds a project tester for testing.
type StructureProviderBuilder func(executable func() (*internal.Executable, error)) internal.ProjectStructureProvider

// InvokeStructure runs the Structure function on a tool and checks the result.
func InvokeStructure(
	t *testing.T,
	builder StructureProviderBuilder,
	info internal.ProjectInfo,
	execMockSetup ExecMockSetup,
	resultCheck func(structure *internal.ProjectStructure, err error),
) internal.ProjectStructureProvider {
	return InvokeToolResult(
		t,
		builder,
		func(tool internal.ProjectStructureProvider) (*internal.ProjectStructure, error) {
			return tool.Structure(t.Context(), info)
		},
		execMockSetup,
		resultCheck,
	)
}
