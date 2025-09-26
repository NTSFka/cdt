package workflow_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"cdt/internal/workflow"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCMakeType_Detect_NoCMakeLists(t *testing.T) {
	workflowType := workflow.CMake{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestCMakeType_Detect_CMakeLists(t *testing.T) {
	workflowType := workflow.CMake{}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0600)
	require.NoError(t, err)

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestCMakeType_Create_CustomBuildDirectory(t *testing.T) {
	workflowType := workflow.CMake{}

	tools := internal.Tools{
		tool.NewCMake(func() *internal.Executable { return &internal.Executable{Path: "cmake-test"} }),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
		tool.NewCTest(func() *internal.Executable { return nil }),
	}

	buildDirectory := "my-build-directory"

	project := workflowType.Create(workflow.Config{Directory: "dir1", IntermediateDirectory: &buildDirectory}, tools)

	assert.Equal(t, "dir1", project.Info.Directory)
	require.NotNil(t, project.Info.IntermediateDirectory)
	assert.Equal(t, buildDirectory, *project.Info.IntermediateDirectory)
	assert.NotNil(t, project.Workflow.Linter)
	assert.NotNil(t, project.Workflow.Formatter)
}

func TestCMakeType_Project_TestAll(t *testing.T) {
	workflowType := workflow.CMake{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0600)
	require.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	project := workflowType.Create(workflow.Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	require.NotNil(t, project.Workflow.Tester)
	cmakeMock.OnRunAnything("cmake-test").Return(nil)
	ctestMock.OnRun("ctest-test", []string{"--test-dir", buildDir}).Return(nil)

	err = project.Workflow.Tester.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: project.Info})
	require.NoError(t, err)

	cmakeMock.AssertExpectations(t)
	ctestMock.AssertExpectations(t)
}

func TestCMakeType_Project_TestAll_BuildFailed(t *testing.T) {
	workflowType := workflow.CMake{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0600)
	require.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := workflowType.Create(workflow.Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	require.NotNil(t, p.Workflow.Tester)
	cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

	err = p.Workflow.Tester.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: p.Info})
	require.EqualError(t, err, "build failed: failed")

	cmakeMock.AssertExpectations(t)
	ctestMock.AssertExpectations(t)
}

func TestCMakeProject_Project_Test(t *testing.T) {
	workflowType := workflow.CMake{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0600)
	require.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	project := workflowType.Create(workflow.Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	require.NotNil(t, project.Workflow.Tester)
	cmakeMock.OnRunAnything("cmake-test").Return(nil)
	ctestMock.OnRun("ctest-test", []string{"-R", "my-test", "--test-dir", buildDir}).Return(nil)

	err = project.Workflow.Tester.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: project.Info},
		"my-test",
	)
	require.NoError(t, err)

	cmakeMock.AssertExpectations(t)
	ctestMock.AssertExpectations(t)
}

func TestCMakeProject_Project_TestBuild_Failed(t *testing.T) {
	workflowType := workflow.CMake{}
	cmakeMock := test.NewExecutable(t)
	ctestMock := test.NewExecutable(t)

	tools := internal.Tools{
		tool.NewCMake(cmakeMock.LazyExecutable("cmake-test")),
		tool.NewCTest(ctestMock.LazyExecutable("ctest-test")),
		tool.NewClangFormat(func() *internal.Executable { return nil }),
		tool.NewClangTidy(func() *internal.Executable { return nil }),
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0600)
	require.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := workflowType.Create(workflow.Config{Directory: dir, IntermediateDirectory: &buildDir}, tools)

	require.NotNil(t, p.Workflow.Tester)
	cmakeMock.OnRunAnything("cmake-test").Return(errors.New("failed"))

	err = p.Workflow.Tester.TestPattern(t.Context(), internal.ProjectTesterOptions{ProjectInfo: p.Info}, "my-test")
	require.EqualError(t, err, "build failed: failed")

	cmakeMock.AssertExpectations(t)
	ctestMock.AssertExpectations(t)
}
