package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runFormatDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "format")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runFormat(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runFormatDir(project, buildDirectory.Path(), args...)
}

// Test format of a project that cannot be formatted
func TestFormatCannotBeFormatted(t *testing.T) {
	err := runFormat(t, "empty", "check")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support source formatting", err.Error())
	}
}

// Test build of a project that uses C++ with CMake
func TestFormatCxxCmake(t *testing.T) {
	checkTool(t, "cmake")
	checkTool(t, "clang-format")

	err := runFormat(t, "cxx-cmake/valid")

	assert.NoError(t, err)
}

// Test build of a project that uses C++ with CMake
func TestFormatFileCxxCmake(t *testing.T) {
	checkTool(t, "cmake")
	checkTool(t, "clang-format")

	err := runFormat(t, "cxx-cmake/valid", "main.cpp")

	assert.NoError(t, err)
}

// Test build of a project that uses C++ with CMake
func TestFormatCheckCxxCmake(t *testing.T) {
	checkTool(t, "cmake")
	checkTool(t, "clang-format")

	err := runFormat(t, "cxx-cmake/valid", "--check")

	assert.NoError(t, err)
}

// Test build of a project that uses C++ with CMake
func TestFormatCheckFileCxxCmake(t *testing.T) {
	checkTool(t, "cmake")
	checkTool(t, "clang-format")

	err := runFormat(t, "cxx-cmake/valid", "--check", "test.cpp")

	assert.NoError(t, err)
}

// Test format of a project that uses C++ with CMake - invalid configuration
func TestFormatCheckCxxCmakeConfigureInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runFormat(t, "cxx-cmake/configure-invalid", "check")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 1", err.Error())
	}
}
