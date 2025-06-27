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
	environment.On("FindExecutable", "cmake").Return(environment.MakeExecutable("/bin/cmake"))

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
	environment.On("FindExecutable", "cmake").Return(nil)

	cmake := DetectCMake(&environment)
	assert.NotNil(t, cmake)
	assert.Equal(t, "cmake", cmake.Id())
	assert.False(t, cmake.IsAvailable())
	assert.Nil(t, cmake.ExecutablePath())

	environment.AssertExpectations(t)
}

func TestCMake_Configure(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	environment.OnRunSuccess(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	})

	err := cmake.Configure(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestCMake_Configure_Failed(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	environment.OnRunError(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}, errors.New("failed"))

	err := cmake.Configure(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMake_Structure_ConfigureFailed(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunError(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}, errors.New("failed"))

	structure, err := cmake.Structure(p)
	assert.Nil(t, structure)
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMake_BuildAll(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunSuccess(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	})

	// Build
	environment.OnRunSuccess(p, "cmake", []string{"--build", p.BuildDirectory()})

	err := cmake.BuildAll(p, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestCMake_BuildAll_Failed(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunSuccess(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	})

	// Build
	environment.OnRunError(p, "cmake", []string{"--build", p.BuildDirectory()}, errors.New("failed"))

	err := cmake.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMake_BuildAll_ConfigureFailed(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunError(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}, errors.New("failed"))

	err := cmake.BuildAll(p, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMake_BuildTargets(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunSuccess(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	})

	// Build
	environment.OnRunSuccess(p, "cmake", []string{
		"--build", p.BuildDirectory(),
		"--target", "target1", "target2",
	})

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.NoError(t, err)

	environment.AssertExpectations(t)
}

func TestCMake_BuildTargets_Failed(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunSuccess(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	})

	// Build
	environment.OnRunError(p, "cmake", []string{
		"--build", p.BuildDirectory(),
		"--target", "target1", "target2",
	}, errors.New("failed"))

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMake_BuildTargets_ConfigureFailed(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunError(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	}, errors.New("failed"))

	err := cmake.BuildTargets(p, []string{"target1", "target2"}, []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}

func TestCMake_RunTarget_Failed(t *testing.T) {
	environment := test.Environment{}

	cmake := NewCMake(environment.MakeExecutable("cmake"))

	p := internal.MakeProject("project", t.TempDir(), cmake, internal.Workflow{})

	// Configure
	environment.OnRunSuccess(p, "cmake", []string{
		"-DCMAKE_EXPORT_COMPILE_COMMANDS=ON",
		"-S", ".",
		"-B", p.BuildDirectory(),
	})

	// Build
	environment.OnRunError(p, "cmake", []string{
		"--build", p.BuildDirectory(),
		"--target", "target1",
	}, errors.New("failed"))

	err := cmake.RunTarget(p, "target1", []string{})
	assert.EqualError(t, err, "failed")

	environment.AssertExpectations(t)
}
