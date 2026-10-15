package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestComposer_DetectComposer_Phar(t *testing.T) {
	composer := test.RunDetectToolLastFound(t, composerDetect, []string{"composer.phar"})

	assert.Equal(t, tool.IdComposer, composer.Id())
}

func TestComposer_DetectComposer_System(t *testing.T) {
	composer := test.RunDetectToolLastFound(
		t,
		composerDetect,
		[]string{"composer.phar", "composer"},
	)

	assert.Equal(t, tool.IdComposer, composer.Id())
}

func TestComposer_DetectComposer_NotFound(t *testing.T) {
	composer := test.RunDetectToolNotFound(t, composerDetect, []string{"composer.phar", "composer"})

	assert.Equal(t, tool.IdComposer, composer.Id())
}

func TestComposer_DetectComposer_Config(t *testing.T) {
	composer := test.RunDetectToolConfig(t, composerDetect, "composer")

	assert.Equal(t, tool.IdComposer, composer.Id())
}

func TestComposer_Composer_AddDependencies(t *testing.T) {
	test.InvokeAddDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		false,
		[]string{"require", "dep1"},
	)
}

func TestComposer_Composer_AddDependencies_Dev(t *testing.T) {
	test.InvokeAddDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		true,
		[]string{"require", "--dev", "dep1"},
	)
}

func TestComposer_Composer_RemoveDependencies(t *testing.T) {
	test.InvokeRemoveDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		false,
		[]string{"remove", "dep1"},
	)
}

func TestComposer_Composer_RemoveDependencies_Dev(t *testing.T) {
	test.InvokeRemoveDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		true,
		[]string{"remove", "--dev", "dep1"},
	)
}

func TestComposer_Composer_UpdateDependencies(t *testing.T) {
	test.InvokeUpdateDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		[]string{"update", "dep1"},
	)
}

func TestComposer_Composer_FetchDependencies(t *testing.T) {
	test.InvokeFetchDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		false,
		[]string{"install"},
	)
}

func TestComposer_Composer_FetchDependencies_NoDev(t *testing.T) {
	test.InvokeFetchDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		true,
		[]string{"install", "--no-dev"},
	)
}

func TestComposer_Composer_ListDependencies(t *testing.T) {
	test.InvokeListDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"show"},
	)
}

func TestComposer_Composer_AuditDependencies(t *testing.T) {
	test.InvokeAuditDependenciesSuccess(
		t,
		composerDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"audit"},
	)
}

func composerDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectComposer(ctx, options)
}

func composerDependencyManager(
	executable func() (*internal.Executable, error),
) internal.ProjectDependencyManager {
	return tool.NewComposer(executable)
}
