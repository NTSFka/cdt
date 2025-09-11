package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCMake_CMakeDetect(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("cmake").
		Return(env.NewExecutable("/bin/cmake"))

	cmake := DetectCMake(context.Background(), env)
	assert.NotNil(t, cmake)
	assert.Equal(t, "cmake", cmake.Id())
	assert.True(t, cmake.IsAvailable())

	if executable := cmake.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/cmake", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestCMake_CMakeDetect_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("cmake").
		Return(nil)

	cmake := DetectCMake(context.Background(), env)
	assert.NotNil(t, cmake)
	assert.Equal(t, "cmake", cmake.Id())
	assert.False(t, cmake.IsAvailable())
	assert.Nil(t, cmake.Executable())

	env.AssertExpectations(t)
}

func TestCMake_Configure(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	info := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).Return(nil)

	err := cmake.Configure(context.Background(), internal.ProjectConfiguratorOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestCMake_Configure_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	info := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).Return(errors.New("failed"))

	err := cmake.Configure(context.Background(), internal.ProjectConfiguratorOptions{ProjectInfo: info})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestCMake_Structure_ConfigureFailed(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(errors.New("failed"))

	structure, err := cmake.Structure(context.Background(), desc)
	assert.Nil(t, structure)
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestCMake_BuildAll(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(nil)

	// Build
	exec.OnRun("cmake", []string{"--build", buildDir}).
		Return(nil)

	err := cmake.BuildAll(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestCMake_BuildAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(nil)

	// Build
	exec.OnRun("cmake", []string{"--build", buildDir}).
		Return(errors.New("failed"))

	err := cmake.BuildAll(context.Background(), desc, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestCMake_BuildAll_ConfigureFailed(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(errors.New("failed"))

	err := cmake.BuildAll(context.Background(), desc, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestCMake_BuildTargets(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(nil)

	// Build
	exec.OnRun("cmake", []string{
		"--build", buildDir,
		"--target", "target1", "target2",
	}).
		Return(nil)

	err := cmake.BuildTargets(context.Background(), desc, []string{"target1", "target2"}, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestCMake_BuildTargets_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(nil)

	// Build
	exec.OnRun("cmake", []string{
		"--build", buildDir,
		"--target", "target1", "target2",
	}).
		Return(errors.New("failed"))

	err := cmake.BuildTargets(context.Background(), desc, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestCMake_BuildTargets_ConfigureFailed(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(errors.New("failed"))

	err := cmake.BuildTargets(context.Background(), desc, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestCMake_RunTarget_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	cmake := NewCMake(exec.LazyExecutable("cmake"))

	buildDir := t.TempDir()

	desc := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: &buildDir,
		StructureProvider:     cmake,
	}

	// Configure
	exec.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", buildDir,
	}).
		Return(nil)

	// Build
	exec.OnRun("cmake", []string{
		"--build", buildDir,
		"--target", "target1",
	}).
		Return(errors.New("failed"))

	err := cmake.RunTarget(context.Background(), desc, "target1", []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
