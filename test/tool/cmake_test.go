package tool

import (
	"cdt/internal"
	"cdt/internal/project"
	"cdt/internal/tool"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"testing"
)

// If a project directory doesn't contain CMakeLists.txt, cmakeTester can't be created
func TestDetectCMakeProjectNoCMakeLists(t *testing.T) {
	tools := internal.Tools{}

	structureProvider, workflow, err := project.DetectCMakeProject("dir1", tools)

	assert.NoError(t, err)
	assert.Nil(t, structureProvider)
	assert.Nil(t, workflow)
}

// If cmake is not available
func TestDetectCMakeProjectNoCMakeBinary(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(nil),
	}

	structureProvider, workflow, err := project.DetectCMakeProject("data/cmake", tools)

	assert.Error(t, err, "cmake is not installed on the system")
	assert.Nil(t, structureProvider)
	assert.Nil(t, workflow)
}

// If no formatters and linters are available
func TestDetectCMakeProjectNoFormatterAndLinters(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test"}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
		&tool.CTest{},
	}

	structureProvider, workflow, err := project.DetectCMakeProject("data/cmake", tools)

	assert.NoError(t, err)
	assert.NotNil(t, structureProvider)
	assert.NotNil(t, workflow)
}

func TestCMakeProjectConfigureAndBuildAndRun(t *testing.T) {
	checkTool(t, "cmake")

	buildDirectory := fs.NewDir(t, "cdt-test")

	cmake := tool.DetectCMake()

	project := internal.MakeProject("data/cmake", buildDirectory.Path(), cmake)

	var err = cmake.Configure(project, []string{})
	assert.NoError(t, err)

	var structure *internal.ProjectStructure
	structure, err = cmake.Structure(project)
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

	err = cmake.BuildAll(project, []string{})
	assert.NoError(t, err)

	err = cmake.BuildTargets(project, []string{"main"}, []string{})
	assert.NoError(t, err)

	err = cmake.RunTarget(project, "main", []string{})
	assert.NoError(t, err)
}
