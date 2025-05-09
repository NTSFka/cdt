package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runBuildDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "build")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runBuild(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runBuildDir(project, buildDirectory.Path(), args...)
}

// Test build of a project that cannot be configured
func TestBuildCannotBeBuilt(t *testing.T) {
	err := runBuild(t, "empty")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support building", err.Error())
	}
}

// Test build of a project that uses C++ with CMake
func TestBuildCxxCmake(t *testing.T) {
	checkTool(t, "cmake")

	err := runBuild(t, "cxx-cmake/valid")
	assert.NoError(t, err)
}

// Test build of a project target that uses C++ with CMake
func TestBuildCxxCmakeTarget(t *testing.T) {
	checkTool(t, "cmake")

	err := runBuild(t, "cxx-cmake/valid", "main")
	assert.NoError(t, err)
}

// Test build of a project that uses C++ with CMake - invalid configuration
func TestBuildCxxCmakeConfigureInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runBuild(t, "cxx-cmake/configure-invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 1", err.Error())
	}
}

// Test build of a project that uses C++ with CMake - invalid build
func TestBuildCxxCmakeInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runBuild(t, "cxx-cmake/build-invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 2", err.Error())
	}
}
