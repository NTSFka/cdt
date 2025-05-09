package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runProjectDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "project")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runProject(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runProjectDir(project, buildDirectory.Path(), args...)
}

func TestProjectTargets(t *testing.T) {
	err := runProject(t, "cxx-cmake/valid", "targets")

	assert.NoError(t, err)
}

func TestProjectFiles(t *testing.T) {
	err := runProject(t, "cxx-cmake/valid", "files")

	assert.NoError(t, err)
}
