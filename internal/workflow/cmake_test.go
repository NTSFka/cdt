package workflow

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCMakeType_Detect_NoCMakeLists(t *testing.T) {
	workflowType := CMake{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestCMakeType_Detect_CMakeLists(t *testing.T) {
	workflowType := CMake{}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestCMakeType_Create_CustomBuildDirectory(t *testing.T) {
	workflowType := CMake{}

	tools := internal.Tools{
		tool.NewCMake(func() *internal.Executable { return &internal.Executable{Path: "cmake-test"} }),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
		tool.NewCTest(func() *internal.Executable { return nil }),
	}

	buildDirectory := "my-build-directory"

	p := workflowType.Create(Config{Directory: "dir1", IntermediateDirectory: &buildDirectory}, tools)

	assert.Equal(t, "dir1", p.Info.Directory)
	if assert.NotNil(t, p.Info.IntermediateDirectory) {
		assert.Equal(t, buildDirectory, *p.Info.IntermediateDirectory)
	}
	assert.NotNil(t, p.Workflow.Linter)
	assert.NotNil(t, p.Workflow.Formatter)
}

func TestCMakeType_Project_TestAll(t *testing.T) {
	workflowType := CMake{}
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

	p := workflowType.Create(Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(nil)
		ctestMock.OnRun("ctest-test", []string{"--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.TestAll(context.Background(), internal.ProjectTesterOptions{ProjectInfo: p.Info})
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMakeType_Project_TestAll_BuildFailed(t *testing.T) {
	workflowType := CMake{}
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

	p := workflowType.Create(Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

		err = p.Workflow.Tester.TestAll(context.Background(), internal.ProjectTesterOptions{ProjectInfo: p.Info})
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMakeProject_Project_Test(t *testing.T) {
	workflowType := CMake{}
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

	p := workflowType.Create(Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(nil)
		ctestMock.OnRun("ctest-test", []string{"-R", "my-test", "--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.TestPattern(context.Background(), internal.ProjectTesterOptions{ProjectInfo: p.Info}, "my-test")
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestCMakeProject_Project_TestBuild_Failed(t *testing.T) {
	workflowType := CMake{}
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

	p := workflowType.Create(Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	if assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

		err = p.Workflow.Tester.TestPattern(context.Background(), internal.ProjectTesterOptions{ProjectInfo: p.Info}, "my-test")
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}
