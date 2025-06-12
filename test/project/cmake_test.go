package project

import (
	"cdt/internal"
	"cdt/internal/project"
	"cdt/internal/tool"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"testing"
)

// If a project directory doesn't contain CMakeLists.txt, CMakeProject can't be created
func TestDetectCMakeProjectNoCMakeLists(t *testing.T) {
	buildDirectory := fs.NewDir(t, "cdt-test")

	tools := internal.Tools{}

	cmakeProject, err := project.DetectCMakeProject("dir1", buildDirectory.Path(), tools)

	assert.NoError(t, err)
	assert.Nil(t, cmakeProject)
}

// If cmake is not available
func TestDetectCMakeProjectNoCMakeBinary(t *testing.T) {
	buildDirectory := fs.NewDir(t, "cdt-test")

	tools := internal.Tools{
		tool.NewCMake(nil),
	}

	cmakeProject, err := project.DetectCMakeProject("data/cmake", buildDirectory.Path(), tools)

	assert.Error(t, err, "cmake is not installed on the system")
	assert.Nil(t, cmakeProject)
}

// If no formatters and linters are available
func TestDetectCMakeProjectNoFormatterAndLinters(t *testing.T) {
	buildDirectory := fs.NewDir(t, "cdt-test")

	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test"}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
	}

	cmakeProject, err := project.DetectCMakeProject("data/cmake", buildDirectory.Path(), tools)

	assert.NoError(t, err)

	if assert.NotNil(t, cmakeProject) {
		assert.Equal(t, cmakeProject.RootDirectory(), "data/cmake")
		assert.Equal(t, cmakeProject.BuildDirectory(), buildDirectory.Path())
	}
}

func TestCMakeProjectConfigureAndBuildAndRun(t *testing.T) {
	checkTool(t, "cmake")

	buildDirectory := fs.NewDir(t, "cdt-test")

	tools := internal.Tools{
		tool.DetectCMake(),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
	}

	var cmakeProject *project.CMakeProject
	var err error
	cmakeProject, err = project.DetectCMakeProject("data/cmake", buildDirectory.Path(), tools)
	assert.NoError(t, err)

	err = cmakeProject.Configure(cmakeProject)
	assert.NoError(t, err)

	var structure *internal.ProjectStructure
	structure, err = cmakeProject.Structure()
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

	err = cmakeProject.BuildAll(cmakeProject)
	assert.NoError(t, err)

	err = cmakeProject.BuildTarget(cmakeProject, "main")
	assert.NoError(t, err)

	err = cmakeProject.Run(cmakeProject, "main", []string{})
	assert.NoError(t, err)
}
