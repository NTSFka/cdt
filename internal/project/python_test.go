package project

import (
	"cdt/internal"
	"cdt/internal/tool"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPythonType_Detect_NoModFile(t *testing.T) {
	projectType := PythonType{}

	res := projectType.Detect("dir1")

	assert.False(t, res)
}

func TestPythonType_Detect_ModFile(t *testing.T) {
	projectType := PythonType{}

	dir := t.TempDir()

	_, err := os.Create(filepath.Join(dir, "pyproject.toml"))
	assert.NoError(t, err)

	res := projectType.Detect(dir)

	assert.True(t, res)
}

func TestPythonType_Create(t *testing.T) {
	projectType := PythonType{}

	tools := internal.Tools{
		tool.NewPython(func() *internal.Executable { return &internal.Executable{Path: "php-test"} }),
		tool.NewPyTest(func() *internal.Executable { return &internal.Executable{Path: "pytest-test"} }),
		tool.NewPip(func() *internal.Executable { return &internal.Executable{Path: "pip-test"} }),
	}

	p := projectType.Create(Config{Directory: "dir1"}, tools)

	if assert.NotNil(t, p) {
		assert.Nil(t, p.Workflow.Configurator)
		assert.Nil(t, p.Workflow.Builder)
		assert.NotNil(t, p.Workflow.Tester)
		assert.Nil(t, p.Workflow.Linter)
		assert.Nil(t, p.Workflow.Formatter)
		assert.NotNil(t, p.Workflow.Runner)
	}
}
