package tool

import (
	"cdt/internal"
	"cdt/internal/tool"
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

	checkTool(t, environment, "cmake")

	buildDirectory := fs.NewDir(t, "cdt-test")

	cmake := tool.DetectCMake(environment)

	desc := internal.Project{
		Directory:             "data/cmake",
		IntermediateDirectory: internal.StrPtr(buildDirectory.Path()),
		StructureProvider:     cmake,
	}

	var err = cmake.Configure(desc, []string{})
	assert.NoError(t, err)

	var structure *internal.ProjectStructure
	structure, err = cmake.Structure(desc)
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

	err = cmake.BuildAll(desc, []string{})
	assert.NoError(t, err)

	err = cmake.BuildTargets(desc, []string{"main"}, []string{})
	assert.NoError(t, err)

	err = cmake.RunTarget(desc, "main", []string{})
	assert.NoError(t, err)
}
