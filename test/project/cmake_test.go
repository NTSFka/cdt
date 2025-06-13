package project

import (
	"cdt/internal"
	"cdt/internal/project"
	"cdt/internal/tool"
	"github.com/stretchr/testify/assert"
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
