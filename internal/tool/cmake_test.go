package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCMake_CMakeDetect(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("cmake").
		Return(environment.MakeExecutable("/bin/cmake"))

	cmake := DetectCMake(&environment)
	assert.NotNil(t, cmake)
	assert.Equal(t, "cmake", cmake.Id())
	assert.True(t, cmake.IsAvailable())

	if path := cmake.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/cmake", *path)
	}

	environment.AssertExpectations(t)
}

func TestCMake_CMakeDetect_NotFound(t *testing.T) {
	environment := test.Environment{}
	environment.OnFindExecutable("cmake").
		Return(nil)

	cmake := DetectCMake(&environment)
	assert.NotNil(t, cmake)
	assert.Equal(t, "cmake", cmake.Id())
	assert.False(t, cmake.IsAvailable())
	assert.Nil(t, cmake.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestCMake_Configure(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).Return(nil)

	err := cmake.Configure(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestCMake_Configure_Failed(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).Return(errors.New("failed"))

	err := cmake.Configure(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestCMake_Structure_ConfigureFailed(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(errors.New("failed"))

	structure, err := cmake.Structure(p)
	assert.Nil(t, structure)
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestCMake_BuildAll(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(nil)

	// Build
	runMock.OnRun("cmake", []string{"--build", p.BuildDirectory()}).
		Return(nil)

	err := cmake.BuildAll(p, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestCMake_BuildAll_Failed(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(nil)

	// Build
	runMock.OnRun("cmake", []string{"--build", p.BuildDirectory()}).
		Return(errors.New("failed"))

	err := cmake.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestCMake_BuildAll_ConfigureFailed(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(errors.New("failed"))

	err := cmake.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestCMake_BuildTargets(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(nil)

	// Build
	runMock.OnRun("cmake", []string{
		"--build", p.BuildDirectory(),
		"--target", "target1", "target2",
	}).
		Return(nil)

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.NoError(t, err)

	runMock.AssertExpectations(t)
}

func TestCMake_BuildTargets_Failed(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(nil)

	// Build
	runMock.OnRun("cmake", []string{
		"--build", p.BuildDirectory(),
		"--target", "target1", "target2",
	}).
		Return(errors.New("failed"))

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestCMake_BuildTargets_ConfigureFailed(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(errors.New("failed"))

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}

func TestCMake_RunTarget_Failed(t *testing.T) {
	runMock := test.Executable{}

	cmake := NewCMake(runMock.LazyExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	runMock.OnRun("cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}).
		Return(nil)

	// Build
	runMock.OnRun("cmake", []string{
		"--build", p.BuildDirectory(),
		"--target", "target1",
	}).
		Return(errors.New("failed"))

	err := cmake.RunTarget(p, "target1", []string{})
	assert.EqualError(t, err, "failed")

	runMock.AssertExpectations(t)
}
