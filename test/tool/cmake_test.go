package tool

import (
	"cdt/internal"
	"cdt/internal/tool"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/fs"
)

func TestCMakeRealProjectConfigureAndBuildAndRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, t.Context(), environment, "cmake")

	buildDirectory := fs.NewDir(t, "cdt-test")

	cmake := tool.DetectCMake(t.Context(), environment)

	info := internal.ProjectInfo{
		Directory:             "data/cmake",
		IntermediateDirectory: internal.StrPtr(buildDirectory.Path()),
		StructureProvider:     cmake,
	}

	var err = cmake.Configure(t.Context(), internal.ProjectConfiguratorOptions{ProjectInfo: info})

	require.NoError(t, err)

	var structure *internal.ProjectStructure
	structure, err = cmake.Structure(t.Context(), info)
	require.NoError(t, err)

	require.NotNil(t, structure)
	assert.Equal(t, map[string]internal.ProjectTarget{
		"fmt": {
			Dependency: true,
			Files:      nil,
		},
		"main": {
			Dependency: false,
			Files:      []string{"main.cpp"},
		},
		"main_test": {
			Dependency: false,
			Files:      []string{"test.cpp"},
		},
	}, structure.Targets)

	assert.Equal(t, []string{"main.cpp", "test.cpp"}, structure.GetFiles())

	err = cmake.BuildAll(t.Context(), internal.ProjectBuilderOptions{ProjectInfo: info})
	require.NoError(t, err)

	err = cmake.BuildTargets(t.Context(), internal.ProjectBuilderOptions{ProjectInfo: info}, []string{"main"})
	require.NoError(t, err)

	err = cmake.RunTarget(t.Context(), internal.ProjectRunnerOptions{ProjectInfo: info}, "main")
	require.NoError(t, err)
}
