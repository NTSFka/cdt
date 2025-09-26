package workflow_test

import (
	"os"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/tool"
	"cdt/internal/workflow"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createGoModFile(dir string) error {
	_, err := os.Create(filepath.Join(dir, "go.mod"))

	return err
}

func TestGoType_Detect_NoModFile(t *testing.T) {
	workflowType := workflow.Go{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestGoType_Detect_ModFile(t *testing.T) {
	workflowType := workflow.Go{}

	dir := t.TempDir()

	err := createGoModFile(dir)
	require.NoError(t, err)

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestGoType_Create(t *testing.T) {
	workflowType := workflow.Go{}

	tools := internal.Tools{
		tool.NewGo(func() *internal.Executable { return &internal.Executable{Path: "go-test"} }),
		tool.NewGolangCILint(func() *internal.Executable { return nil }),
	}

	project := workflowType.Create(workflow.Config{Directory: "dir1"}, tools)

	require.NotNil(t, project)
	assert.Nil(t, project.Workflow.Configurator)
	assert.NotNil(t, project.Workflow.Builder)
	assert.NotNil(t, project.Workflow.Tester)
	assert.NotNil(t, project.Workflow.Linter)
	assert.NotNil(t, project.Workflow.Formatter)
	assert.NotNil(t, project.Workflow.Runner)
	assert.NotNil(t, project.Workflow.DependencyManager)
}
