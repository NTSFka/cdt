package tool

import (
	"cdt/internal"
	"cdt/internal/tool"
	"context"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
)

func TestCMakeRealProjectConfigureAndBuildAndRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, context.Background(), environment, "cmake")

	buildDirectory := fs.NewDir(t, "cdt-test")

	cmake := tool.DetectCMake(context.Background(), environment)

	info := internal.ProjectInfo{
		Directory:             "data/cmake",
		IntermediateDirectory: internal.StrPtr(buildDirectory.Path()),
		StructureProvider:     cmake,
	}

	var err = cmake.Configure(context.Background(), internal.ProjectConfiguratorOptions{ProjectInfo: info})
	assert.NoError(t, err)

	var structure *internal.ProjectStructure
	structure, err = cmake.Structure(context.Background(), info)
	assert.NoError(t, err)

	if assert.NotNil(t, structure) {
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
	}

	err = cmake.BuildAll(context.Background(), info, []string{})
	assert.NoError(t, err)

	err = cmake.BuildTargets(context.Background(), info, []string{"main"}, []string{})
	assert.NoError(t, err)

	err = cmake.RunTarget(context.Background(), info, "main", []string{})
	assert.NoError(t, err)
}
