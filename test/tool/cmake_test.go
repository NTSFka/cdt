package tool

import (
	"cdt/internal"
	"cdt/internal/tool"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"testing"
)

func TestCMakeProjectConfigureAndBuildAndRun(t *testing.T) {
	checkTool(t, "cmake")

	buildDirectory := fs.NewDir(t, "cdt-test")

	cmake := tool.DetectCMake()

	p := internal.MakeProject("data/cmake", buildDirectory.Path(), cmake, internal.Workflow{})

	var err = cmake.Configure(p, []string{})
	assert.NoError(t, err)

	var structure *internal.ProjectStructure
	structure, err = cmake.Structure(p)
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

	err = cmake.BuildAll(p, []string{})
	assert.NoError(t, err)

	err = cmake.BuildTargets(p, []string{"main"}, []string{})
	assert.NoError(t, err)

	err = cmake.RunTarget(p, "main", []string{})
	assert.NoError(t, err)
}
