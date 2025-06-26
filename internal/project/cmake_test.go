package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// If a project directory doesn't contain CMakeLists.txt, cmakeTester can't be created
func TestDetectCMakeProjectNoCMakeLists(t *testing.T) {
	tools := internal.Tools{}

	p := DetectCMakeProject(internal.Config{RootDirectory: "dir1"}, tools)

	assert.Nil(t, p)
}

// If cmake is not available
func TestDetectCMakeProjectNoCMakeBinary(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(nil),
	}

	p := DetectCMakeProject(internal.Config{RootDirectory: "data/cmake"}, tools)

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

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	p := DetectCMakeProject(internal.Config{RootDirectory: dir}, tools)

	if assert.NotNil(t, p) {
		assert.Nil(t, p.Workflow.Linter)
		assert.Nil(t, p.Workflow.Formatter)
	}
}

// If formatters and linters are available
func TestDetectCMakeProjectFormatterAndLinters(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test"}),
		tool.NewClangFormat(&internal.Executable{Path: "clang-format-test"}, nil),
		tool.NewClangTidy(&internal.Executable{Path: "clang-tidy-test"}, nil),
		&tool.CTest{},
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	p := DetectCMakeProject(internal.Config{RootDirectory: dir}, tools)

	if assert.NotNil(t, p) {
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
	}
}
