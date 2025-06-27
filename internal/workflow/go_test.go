package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func createGoModFile(dir string) error {
	_, err := os.Create(filepath.Join(dir, "go.mod"))

	return err
}

func TestGo_DetectGoProject_NoModFile(t *testing.T) {
	tools := internal.Tools{}

	p := DetectGoProject(internal.Config{RootDirectory: "dir1"}, tools)

	assert.Nil(t, p)
}

func TestGo_DetectGpProject_NoGoBinary(t *testing.T) {
	tools := internal.Tools{
		tool.NewGo(nil),
	}

	dir := t.TempDir()

	err := createGoModFile(dir)
	assert.NoError(t, err)

	p := DetectGoProject(internal.Config{RootDirectory: dir}, tools)

	assert.Nil(t, p)
}

func TestGo_DetectGoProject_NoLinter(t *testing.T) {
	tools := internal.Tools{
		tool.NewGo(&internal.Executable{Path: "go-test"}),
		&tool.GolangCILint{},
	}

	dir := t.TempDir()

	err := createGoModFile(dir)
	assert.NoError(t, err)

	assert.True(t, internal.PathExists(filepath.Join(dir, "go.mod")))
	p := DetectGoProject(internal.Config{RootDirectory: dir}, tools)

	if assert.NotNil(t, p) {
		assert.Nil(t, p.Workflow.Configurator)
		assert.NotNil(t, p.Workflow.Builder)
		assert.NotNil(t, p.Workflow.Tester)
		assert.Nil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
		assert.NotNil(t, p.Workflow.Runner)
	}
}

func TestGo_DetectGoProject_Linter(t *testing.T) {
	tools := internal.Tools{
		tool.NewGo(&internal.Executable{Path: "go-test"}),
		tool.NewGolangCILint(&internal.Executable{Path: "golangci-lint-test"}),
	}

	dir := t.TempDir()

	err := createGoModFile(dir)
	assert.NoError(t, err)

	assert.True(t, internal.PathExists(filepath.Join(dir, "go.mod")))
	p := DetectGoProject(internal.Config{RootDirectory: dir}, tools)

	if assert.NotNil(t, p) {
		assert.Nil(t, p.Workflow.Configurator)
		assert.NotNil(t, p.Workflow.Builder)
		assert.NotNil(t, p.Workflow.Tester)
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
		assert.NotNil(t, p.Workflow.Runner)
	}
}
