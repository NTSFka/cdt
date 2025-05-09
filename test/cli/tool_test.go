package cli

import (
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func runToolDir(project string, buildDirectory string, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "--directory", filepath.Join("data", project))
	runArgs = append(runArgs, "--build", buildDirectory)
	runArgs = append(runArgs, "tool")
	runArgs = append(runArgs, args...)

	return runMain(runArgs...)
}

func runTool(t *testing.T, project string, args ...string) error {
	buildDirectory := fs.NewDir(t, filepath.Join("cdt-test", project))

	return runToolDir(project, buildDirectory.Path(), args...)
}

func TestToolList(t *testing.T) {
	err := runTool(t, "cxx-cmake/valid", "list")

	assert.NoError(t, err)
}

func TestToolListAll(t *testing.T) {
	err := runTool(t, "cxx-cmake/valid", "list", "--all")

	assert.NoError(t, err)
}
