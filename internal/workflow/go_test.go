package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createGoModFile(dir string) error {
	_, err := os.Create(filepath.Join(dir, "go.mod"))

	return err
}

func TestGoType_Detect_NoModFile(t *testing.T) {
	workflowType := Go{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestGoType_Detect_ModFile(t *testing.T) {
	workflowType := Go{}

	dir := t.TempDir()

	err := createGoModFile(dir)
	assert.NoError(t, err)

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestGoType_Create(t *testing.T) {
	workflowType := Go{}

	tools := internal.Tools{
		tool.NewGo(func() *internal.Executable { return &internal.Executable{Path: "go-test"} }),
		tool.NewGolangCILint(func() *internal.Executable { return nil }),
	}

	p := workflowType.Create(Config{Directory: "dir1"}, tools)

	if assert.NotNil(t, p) {
		assert.Nil(t, p.Workflow.Configurator)
		assert.NotNil(t, p.Workflow.Builder)
		assert.NotNil(t, p.Workflow.Tester)
		assert.NotNil(t, p.Workflow.Linter)
		assert.NotNil(t, p.Workflow.Formatter)
		assert.NotNil(t, p.Workflow.Runner)
		assert.NotNil(t, p.Workflow.DependencyManager)
	}
}
