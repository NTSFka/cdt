package project

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCMakeType_Detect_NoCMakeLists(t *testing.T) {
	projectType := CMakeType{}

	res := projectType.Detect("dir1")

	assert.False(t, res)
}

func TestCMakeType_Detect_CMakeLists(t *testing.T) {
	projectType := CMakeType{}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	res := projectType.Detect(dir)

	assert.True(t, res)
}

func TestCMakeType_Create_CustomBuildDirectory(t *testing.T) {
	projectType := CMakeType{}

	tools := internal.Tools{
		tool.NewCMake(func() *internal.Executable { return &internal.Executable{Path: "cmake-test"} }),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
		tool.NewCTest(func() *internal.Executable { return nil }),
	}

	buildDirectory := "my-build-directory"

	p := projectType.Create(Config{Directory: "dir1", BuildDirectory: &buildDirectory}, tools)

	assert.Equal(t, "dir1", p.Info.Directory)
	if assert.NotNil(t, p.Info.IntermediateDirectory) {
		assert.Equal(t, buildDirectory, *p.Info.IntermediateDirectory)
	}
	assert.NotNil(t, p.Workflow.Linter)
	assert.NotNil(t, p.Workflow.Formatter)
}

func TestCMakeType_Project_TestAll(t *testing.T) {
	projectType := CMakeType{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

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

	p := projectType.Create(Config{Directory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(nil)
		ctestMock.OnRun("ctest-test", []string{"--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.TestAll(p.Info, []string{})
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMakeType_Project_TestAll_BuildFailed(t *testing.T) {
	projectType := CMakeType{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

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

	p := projectType.Create(Config{Directory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

		err = p.Workflow.Tester.TestAll(p.Info, []string{})
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMakeProject_Project_Test(t *testing.T) {
	projectType := CMakeType{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

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

	p := projectType.Create(Config{Directory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(nil)
		ctestMock.OnRun("ctest-test", []string{"-R", "my-test", "--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.Test(p.Info, "my-test", []string{})
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMakeProject_Project_TestBuild_Failed(t *testing.T) {
	projectType := CMakeType{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

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

	p := projectType.Create(Config{Directory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

		err = p.Workflow.Tester.Test(p.Info, "my-test", []string{})
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}
