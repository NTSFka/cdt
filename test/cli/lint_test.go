package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runLintDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "lint")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runLint(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runLintDir(project, buildDirectory.Path(), args...)
}

// Test lint of a project that cannot be linted
func TestLintCannotBeLinted(t *testing.T) {
	err := runLint(t, "empty")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support linting", err.Error())
	}
}

// Test lint of a project that uses C++ with CMake
func TestLintCxxCmake(t *testing.T) {
	checkTool(t, "cmake")
	checkTool(t, "clang-tidy")

	err := runLint(t, "cxx-cmake/valid")

	assert.NoError(t, err)
}

// Test lint of a project that uses C++ with CMake
func TestLintFileCxxCmake(t *testing.T) {
	checkTool(t, "cmake")
	checkTool(t, "clang-tidy")

	err := runLint(t, "cxx-cmake/valid", "main.cpp")

	assert.NoError(t, err)
}

// Test lint of a project that uses C++ with CMake - invalid configuration
func TestLintCxxCmakeConfigureInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runLint(t, "cxx-cmake/configure-invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 1", err.Error())
	}
}
