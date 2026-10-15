package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPip_DetectPip_Pip(t *testing.T) {
	pip := test.RunDetectToolLastFound(t, pipDetect, []string{"pip"})

	assert.Equal(t, tool.IdPip, pip.Id())
}

func TestPip_DetectPip_System(t *testing.T) {
	pip := test.RunDetectToolLastFound(t, pipDetect, []string{"pip", "pip3"})

	assert.Equal(t, tool.IdPip, pip.Id())
}

func TestPip_DetectPip_NotFound(t *testing.T) {
	pip := test.RunDetectToolNotFound(t, pipDetect, []string{"pip", "pip3"})

	assert.Equal(t, tool.IdPip, pip.Id())
}

func TestPip_DetectPip_Config(t *testing.T) {
	pip := test.RunDetectToolConfig(t, pipDetect, "pip")

	assert.Equal(t, tool.IdPip, pip.Id())
}

func TestPip_Pip_AddDependencies(t *testing.T) {
	test.InvokeAddDependenciesSuccess(
		t,
		pipDependencyManager,
		internal.ProjectDependencyManagerOptions{},
		[]string{"dep1"},
		false,
		[]string{"install", "dep1"},
	)
}

func TestPip_Pip_RemoveDependencies(t *testing.T) {
	test.InvokeRemoveDependenciesSuccess(
		t,
		pipDependencyManager,
		internal.ProjectDependencyManagerOptions{},
		[]string{"dep1"},
		false,
		[]string{"uninstall", "dep1"},
	)
}

func TestPip_Pip_UpdateDependencies(t *testing.T) {
	test.InvokeUpdateDependenciesSuccess(
		t,
		pipDependencyManager,
		internal.ProjectDependencyManagerOptions{},
		[]string{"dep1"},
		[]string{"install", "--upgrade", "dep1"},
	)
}

func TestPip_Pip_FetchDependencies(t *testing.T) {
	test.InvokeFetchDependenciesSuccess(
		t,
		pipDependencyManager,
		internal.ProjectDependencyManagerOptions{},
		false,
		[]string{"install", "-r", "requirements.txt"},
	)
}

func TestPip_Pip_ListDependencies(t *testing.T) {
	test.InvokeListDependenciesSuccess(
		t,
		pipDependencyManager,
		internal.ProjectDependencyManagerOptions{},
		[]string{"list"},
	)
}

func TestPip_Pip_AuditDependencies(t *testing.T) {
	test.InvokeAuditDependenciesError(
		t,
		pipDependencyManager,
		internal.ProjectDependencyManagerOptions{},
		func(err error) {
			require.EqualError(t, err, "not supported")
		},
	)
}

func pipDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPip(ctx, options)
}

func pipDependencyManager(
	executable func() (*internal.Executable, error),
) internal.ProjectDependencyManager {
	return tool.NewPip(executable)
}
