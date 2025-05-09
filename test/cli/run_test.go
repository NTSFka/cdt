package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runRunDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "run")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runRun(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runRunDir(project, buildDirectory.Path(), args...)
}

func TestRunCannotBeBuilt(t *testing.T) {
	err := runRun(t, "empty")

	if assert.Error(t, err) {
		assert.Equal(t, "project doesn't support run of target", err.Error())
	}
}

func TestRunTargetCxxCmakeMissing(t *testing.T) {
	checkTool(t, "cmake")

	err := runRun(t, "cxx-cmake/valid")

	if assert.Error(t, err) {
		assert.Equal(t, "target is required", err.Error())
	}
}

func TestRunTargetCxxCmake(t *testing.T) {
	checkTool(t, "cmake")

	err := runRun(t, "cxx-cmake/valid", "main")
	assert.NoError(t, err)
}

func TestRunCxxCmakeConfigureInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runRun(t, "cxx-cmake/configure-invalid", "main")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 1", err.Error())
	}
}

func TestRunCxxCmakeInvalid(t *testing.T) {
	checkTool(t, "cmake")

	err := runRun(t, "cxx-cmake/build-invalid", "main")

	if assert.Error(t, err) {
		assert.Equal(t, "exit status 2", err.Error())
	}
}
