package workflow

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// If a project directory doesn't contain CMakeLists.txt, cmakeTester can't be created
func TestCMake_DetectCMakeProject_NoCMakeLists(t *testing.T) {
	tools := internal.Tools{}

	p := DetectCMakeProject(internal.Config{RootDirectory: "dir1"}, tools)

	assert.Nil(t, p)
}

func TestCMake_DetectCMakeProject_CustomBuildDirectory(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(func() *internal.Executable { return &internal.Executable{Path: "cmake-test"} }),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
		tool.NewCTest(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDirectory := "my-build-directory"

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDirectory}, tools)

	if assert.NotNil(t, p) {
		assert.Equal(t, buildDirectory, p.BuildDirectory())
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
	}
}

// If no formatters and linters are available
func TestCMake_DetectCMakeProject_NoFormatterAndLinters(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(func() *internal.Executable { return &internal.Executable{Path: "cmake-test"} }),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
		tool.NewCTest(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	p := DetectCMakeProject(internal.Config{RootDirectory: dir}, tools)

	if assert.NotNil(t, p) {
		assert.Equal(t, filepath.Join("build", "dev"), p.BuildDirectory())
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
	}
}

// If formatters and linters are available
func TestCMake_DetectCMakeProject_FormatterAndLinters(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(func() *internal.Executable { return &internal.Executable{Path: "cmake-test"} }),
		tool.NewClangFormat(func() *internal.Executable { return &internal.Executable{Path: "clang-format-test"} }),
		tool.NewClangTidy(func() *internal.Executable { return &internal.Executable{Path: "clang-tidy-test"} }),
		tool.NewCTest(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	p := DetectCMakeProject(internal.Config{RootDirectory: dir}, tools)

	if assert.NotNil(t, p) {
		assert.Equal(t, filepath.Join("build", "dev"), p.BuildDirectory())
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
	}
}

func TestCMake_DetectCMakeProject_TestAll(t *testing.T) {
	var cmakeMock test.Executable
	var ctestMock test.Executable

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(nil)
		ctestMock.OnRun("ctest-test", []string{"--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.TestAll(*p, []string{})
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMake_DetectCMakeProject_TestAll_BuildFailed(t *testing.T) {
	var cmakeMock test.Executable
	var ctestMock test.Executable

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

		err = p.Workflow.Tester.TestAll(*p, []string{})
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMake_DetectCMakeProject_Test(t *testing.T) {
	var cmakeMock test.Executable
	var ctestMock test.Executable

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(nil)
		ctestMock.OnRun("ctest-test", []string{"-R", "my-test", "--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.Test(*p, "my-test", []string{})
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMake_DetectCMakeProject_TestBuild_Failed(t *testing.T) {
	var cmakeMock test.Executable
	var ctestMock test.Executable

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

		err = p.Workflow.Tester.Test(*p, "my-test", []string{})
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}
