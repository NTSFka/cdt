package tool

import (
	"cdt/internal"
	"cdt/internal/tool"
	"errors"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"runtime"
	"testing"
)

func TestCMakeDetect(t *testing.T) {
	environment := testEnvironment{}
	environment.On("FindExecutable", "cmake").Return(environment.makeExecutable("/bin/cmake"))

	cmake := tool.DetectCMake(&environment)
	assert.NotNil(t, cmake)
	assert.Equal(t, "cmake", cmake.Id())
	assert.True(t, cmake.IsAvailable())

	if path := cmake.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/cmake", *path)
	}

	environment.AssertExpectations(t)
}

func TestCMakeDetectNotFound(t *testing.T) {
	environment := testEnvironment{}
	environment.On("FindExecutable", "cmake").Return(nil)

	cmake := tool.DetectCMake(&environment)
	assert.NotNil(t, cmake)
	assert.Equal(t, "cmake", cmake.Id())
	assert.False(t, cmake.IsAvailable())
	assert.Nil(t, cmake.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestCMakeConfigure(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(nil)

	err := cmake.Configure(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestCMakeConfigureFailed(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(errors.New("failed"))

	err := cmake.Configure(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMakeStructureConfigureFailed(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(errors.New("failed"))

	structure, err := cmake.Structure(p)
	assert.Nil(t, structure)
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMakeBuildAll(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(nil)

	// Build
	environment.On("RunExecutable", "cmake", []string{
		"--build", buildDirectory,
	}).Return(nil)

	err := cmake.BuildAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestCMakeBuildAllFailed(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(nil)

	// Build
	environment.On("RunExecutable", "cmake", []string{
		"--build", buildDirectory,
	}).Return(errors.New("failed"))

	err := cmake.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMakeBuildAllConfigureFailed(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(errors.New("failed"))

	err := cmake.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMakeBuildTargets(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(nil)

	// Build
	environment.On("RunExecutable", "cmake", []string{
		"--build", buildDirectory,
		"--target", "target1", "target2",
	}).Return(nil)

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestCMakeBuildTargetsFailed(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(nil)

	// Build
	environment.On("RunExecutable", "cmake", []string{
		"--build", buildDirectory,
		"--target", "target1", "target2",
	}).Return(errors.New("failed"))

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMakeBuildTargetsConfigureFailed(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(errors.New("failed"))

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMakeRunTargetFailed(t *testing.T) {
	environment := testEnvironment{}

	cmake := tool.NewCMake(environment.makeExecutable("cmake"))

	rootDirectory := "project"
	buildDirectory := fs.NewDir(t, "cdt-test").Path()

	p := internal.MakeProject(rootDirectory, buildDirectory, cmake, internal.Workflow{})

	// Configure
	environment.On("RunExecutable", "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", rootDirectory,
		"-B", buildDirectory,
	}).Return(nil)

	// Build
	environment.On("RunExecutable", "cmake", []string{
		"--build", buildDirectory,
		"--target", "target1",
	}).Return(errors.New("failed"))

	err := cmake.RunTarget(p, "target1", []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMakeRealProjectConfigureAndBuildAndRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("project structure is not detected properly on Windows using MSVC Generator")
	}

	environment := internal.SystemEnvironment

	checkTool(t, environment, "cmake")

	buildDirectory := fs.NewDir(t, "cdt-test")

	cmake := tool.DetectCMake(environment)

	project := internal.MakeProject("data/cmake", buildDirectory.Path(), cmake, internal.Workflow{})

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
