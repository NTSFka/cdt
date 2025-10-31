package workflow_test

import (
	"cdt/internal/test"
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
	file, err := os.Create(filepath.Join(dir, "go.mod"))

	if err != nil {
		return err
	}

	return file.Close()
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
		tool.NewGo(test.LazyExecutable("go-test")),
		tool.NewGolangCILint(test.LazyExecutableNil),
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
