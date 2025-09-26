package workflow

import (
	"cdt/internal"
	"cdt/internal/tool"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPythonType_Detect_NoModFile(t *testing.T) {
	workflowType := Python{}

	res := workflowType.Detect("dir1")

	assert.False(t, res)
}

func TestPythonType_Detect_ModFile(t *testing.T) {
	workflowType := Python{}

	dir := t.TempDir()

	_, err := os.Create(filepath.Join(dir, "pyproject.toml"))
	require.NoError(t, err)

	res := workflowType.Detect(dir)

	assert.True(t, res)
}

func TestPythonType_Create(t *testing.T) {
	workflowType := Python{}

	tools := internal.Tools{
		tool.NewPython(func() *internal.Executable { return &internal.Executable{Path: "php-test"} }),
		tool.NewPyTest(func() *internal.Executable { return &internal.Executable{Path: "pytest-test"} }),
		tool.NewPip(func() *internal.Executable { return &internal.Executable{Path: "pip-test"} }),
		tool.NewPylint(func() *internal.Executable { return &internal.Executable{Path: "pylint-test"} }),
		tool.NewFlake8(func() *internal.Executable { return &internal.Executable{Path: "flake8-test"} }),
		tool.NewMyPy(func() *internal.Executable { return &internal.Executable{Path: "mypy-test"} }),
		tool.NewRuff(func() *internal.Executable { return &internal.Executable{Path: "ruff-test"} }),
		tool.NewBandit(func() *internal.Executable { return &internal.Executable{Path: "bandit-test"} }),
		tool.NewBlack(func() *internal.Executable { return &internal.Executable{Path: "black-test"} }),
	}

	project := workflowType.Create(Config{Directory: "dir1"}, tools)

	require.NotNil(t, project)
	assert.Nil(t, project.Workflow.Configurator)
	assert.Nil(t, project.Workflow.Builder)
	assert.NotNil(t, project.Workflow.Tester)
	assert.NotNil(t, project.Workflow.Linter)
	assert.NotNil(t, project.Workflow.Formatter)
	assert.NotNil(t, project.Workflow.Runner)
}
