package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	p := DetectCMakeProject(internal.Config{RootDirectory: dir}, tools)

	assert.Nil(t, p)
}

func TestDetectCMakeProjectCustomBuildDirectory(t *testing.T) {
	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test"}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
		&tool.CTest{},
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDirectory := "my-build-directory"

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDirectory}, tools)

	if assert.NotNil(t, p) {
		assert.Equal(t, buildDirectory, p.BuildDirectory())
		assert.Nil(t, p.Workflow.Linter)
		assert.Nil(t, p.Workflow.Formatter)
	}
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
		assert.Equal(t, filepath.Join("build", "dev"), p.BuildDirectory())
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
		assert.Equal(t, filepath.Join("build", "dev"), p.BuildDirectory())
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
	}
}

func TestDetectCMakeProjectTestAll(t *testing.T) {
	var cmakeMock mock.Mock
	var ctestMock mock.Mock

	cmakeRun := func(ctx internal.RunContext, path string, args []string) error {
		return cmakeMock.Called(ctx, path, args).Error(0)
	}

	ctestRun := func(ctx internal.RunContext, path string, args []string) error {
		return ctestMock.Called(ctx, path, args).Error(0)
	}

	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test", RunFunc: cmakeRun}),
		tool.NewCTest(&internal.Executable{Path: "ctest-test", RunFunc: ctestRun}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.On("func1", mock.Anything, "cmake-test", mock.Anything).Return(nil)
		ctestMock.On("func2", mock.Anything, "ctest-test", []string{"--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.TestAll(*p, []string{})
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestDetectCMakeProjectTestAllBuildFailed(t *testing.T) {
	var cmakeMock mock.Mock
	var ctestMock mock.Mock

	cmakeRun := func(ctx internal.RunContext, path string, args []string) error {
		return cmakeMock.Called(ctx, path, args).Error(0)
	}

	ctestRun := func(ctx internal.RunContext, path string, args []string) error {
		return ctestMock.Called(ctx, path, args).Error(0)
	}

	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test", RunFunc: cmakeRun}),
		tool.NewCTest(&internal.Executable{Path: "ctest-test", RunFunc: ctestRun}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.On("func1", mock.Anything, "cmake-test", mock.Anything).Return(errors.New("failed"))

		err = p.Workflow.Tester.TestAll(*p, []string{})
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestDetectCMakeProjectTest(t *testing.T) {
	var cmakeMock mock.Mock
	var ctestMock mock.Mock

	cmakeRun := func(ctx internal.RunContext, path string, args []string) error {
		return cmakeMock.Called(ctx, path, args).Error(0)
	}

	ctestRun := func(ctx internal.RunContext, path string, args []string) error {
		return ctestMock.Called(ctx, path, args).Error(0)
	}

	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test", RunFunc: cmakeRun}),
		tool.NewCTest(&internal.Executable{Path: "ctest-test", RunFunc: ctestRun}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.On("func1", mock.Anything, "cmake-test", mock.Anything).Return(nil)
		ctestMock.On("func2", mock.Anything, "ctest-test", []string{"-R", "my-test", "--test-dir", buildDir}).Return(nil)

		err = p.Workflow.Tester.Test(*p, "my-test", []string{})
		assert.NoError(t, err)

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}

func TestDetectCMakeProjectTestBuildFailed(t *testing.T) {
	var cmakeMock mock.Mock
	var ctestMock mock.Mock

	cmakeRun := func(ctx internal.RunContext, path string, args []string) error {
		return cmakeMock.Called(ctx, path, args).Error(0)
	}

	ctestRun := func(ctx internal.RunContext, path string, args []string) error {
		return ctestMock.Called(ctx, path, args).Error(0)
	}

	tools := internal.Tools{
		tool.NewCMake(&internal.Executable{Path: "cmake-test", RunFunc: cmakeRun}),
		tool.NewCTest(&internal.Executable{Path: "ctest-test", RunFunc: ctestRun}),
		&tool.ClangFormat{},
		&tool.ClangTidy{},
	}

	dir := t.TempDir()

	_, err := os.OpenFile(filepath.Join(dir, "CMakeLists.txt"), os.O_RDONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	buildDir := filepath.Join(dir, "build")

	p := DetectCMakeProject(internal.Config{RootDirectory: dir, BuildDirectory: &buildDir}, tools)

	if assert.NotNil(t, p) && assert.NotNil(t, p.Workflow.Tester) {
		cmakeMock.On("func1", mock.Anything, "cmake-test", mock.Anything).Return(errors.New("failed"))

		err = p.Workflow.Tester.Test(*p, "my-test", []string{})
		assert.EqualError(t, err, "build failed: failed")

		cmakeMock.AssertExpectations(t)
		ctestMock.AssertExpectations(t)
	}
}
