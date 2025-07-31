package project

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

func TestGoType_Detect_NoModFile(t *testing.T) {
	projectType := GoType{}

	res := projectType.Detect("dir1")

	assert.False(t, res)
}

func TestGoType_Detect_ModFile(t *testing.T) {
	projectType := GoType{}

	dir := t.TempDir()

	err := createGoModFile(dir)
	assert.NoError(t, err)

	res := projectType.Detect(dir)

	assert.True(t, res)
}

func TestGoType_Create(t *testing.T) {
	projectType := GoType{}

	tools := internal.Tools{
		tool.NewGo(func() *internal.Executable { return &internal.Executable{Path: "go-test"} }),
		tool.NewGolangCILint(func() *internal.Executable { return nil }),
	}

	p := projectType.Create(Config{Directory: "dir1"}, tools)

	if assert.NotNil(t, p) {
		assert.Nil(t, p.Workflow.Configurator)
		assert.NotNil(t, p.Workflow.Builder)
		assert.NotNil(t, p.Workflow.Tester)
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
		assert.NotNil(t, p.Workflow.Runner)
	}
}
