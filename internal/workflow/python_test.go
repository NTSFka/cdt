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

func TestPythonType_Detect_NoModFile(t *testing.T) {
	workflowType := workflow.Python{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestPythonType_Detect_ModFile(t *testing.T) {
	workflowType := workflow.Python{}

	dir := t.TempDir()

	_, err := os.Create(filepath.Join(dir, "pyproject.toml"))
	require.NoError(t, err)

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestPythonType_Create(t *testing.T) {
	workflowType := workflow.Python{}

	tools := internal.Tools{
		tool.NewPython(test.LazyExecutable("php-test")),
		tool.NewPyTest(test.LazyExecutable("pytest-test")),
		tool.NewPip(test.LazyExecutable("pip-test")),
		tool.NewPylint(test.LazyExecutable("pylint-test")),
		tool.NewFlake8(test.LazyExecutable("flake8-test")),
		tool.NewMyPy(test.LazyExecutable("mypy-test")),
		tool.NewRuff(test.LazyExecutable("ruff-test")),
		tool.NewBandit(test.LazyExecutable("bandit-test")),
		tool.NewBlack(test.LazyExecutable("black-test")),
	}

	project := workflowType.Create(workflow.Config{Directory: "dir1"}, tools)

	require.NotNil(t, project)
	assert.Nil(t, project.Workflow.Configurator)
	assert.Nil(t, project.Workflow.Builder)
	assert.NotNil(t, project.Workflow.Tester)
	assert.NotNil(t, project.Workflow.Linter)
	assert.NotNil(t, project.Workflow.Formatter)
	assert.NotNil(t, project.Workflow.Runner)
}
