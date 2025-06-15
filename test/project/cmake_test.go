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

	p := project.DetectCMakeProject(internal.Config{RootDirectory: "dir1"}, tools)

	assert.Nil(t, p)
}

// If cmake is not available
func TestDetectCMakeProjectNoCMakeBinary(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(nil),
	}

	p := project.DetectCMakeProject(internal.Config{RootDirectory: "data/cmake"}, tools)

	assert.Nil(t, p)
}

// If no formatters and linters are available
func TestDetectCMakeProjectNoFormatterAndLinters(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test"}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
		&tool.CTest{},
	}

	p := project.DetectCMakeProject(internal.Config{RootDirectory: "data/cmake"}, tools)

	assert.NotNil(t, p)
	assert.Nil(t, p.Workflow.Linter)
	assert.Nil(t, p.Workflow.Formatter)
}

// If formatters and linters are available
func TestDetectCMakeProjectFormatterAndLinters(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test"}),
		tool.NewClangFormat(&internal.Executable{Path: "clang-format-test"}, nil),
		tool.NewClangTidy(&internal.Executable{Path: "clang-tidy-test"}, nil),
		&tool.CTest{},
	}

	p := project.DetectCMakeProject(internal.Config{RootDirectory: "data/cmake"}, tools)

	assert.NotNil(t, p)
	assert.NotNil(t, p.Workflow.Linter)
	assert.NotNil(t, p.Workflow.Formatter)
}
