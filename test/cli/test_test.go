package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runTestDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "test")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runTest(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runTestDir(project, buildDirectory.Path(), args...)
}

// Test "test" of a project that cannot be configured
func TestTestCannotBeBuilt(t *testing.T) {
	err := runTest(t, "empty")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support testing", err.Error())
	}
}

// Test "test" of a project that uses C++ with CMake
func TestTestCxxCmake(t *testing.T) {
	checkTool(t, "cmake")

	err := runTest(t, "cxx-cmake/valid")
	assert.NoError(t, err)
}

// Test "test" of a project that uses C++ with CMake - pattern
func TestTestCxxCmakePattern(t *testing.T) {
	checkTool(t, "cmake")

	err := runTest(t, "cxx-cmake/valid", "test_test")
	assert.NoError(t, err)
}

// Test "test" of a project that uses C++ with CMake - invalid configuration
func TestTestCxxCmakeConfigureInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runTest(t, "cxx-cmake/configure-invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 1", err.Error())
	}
}

// Test "test" of a project that uses C++ with CMake - invalid build
func TestTestCxxCmakeInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runTest(t, "cxx-cmake/build-invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 2", err.Error())
	}
}
