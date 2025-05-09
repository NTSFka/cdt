package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runConfigureDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "configure")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runConfigure(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runConfigureDir(project, buildDirectory.Path(), args...)
}

// Test configuration of a project that cannot be configured
func TestConfigureCannotBeConfigured(t *testing.T) {
	err := runConfigure(t, "empty")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support configuration", err.Error())
	}
}

// Test configuration of a project that uses C++ with CMake
func TestConfigureCxxCmake(t *testing.T) {
	checkTool(t, "cmake")

	err := runConfigure(t, "cxx-cmake/valid")

	assert.NoError(t, err)
}

// Test configuration of a project that uses C++ with CMake - invalid project
func TestConfigureCxxCmakeInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runConfigure(t, "cxx-cmake/configure-invalid")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 1", err.Error())
	}
}
