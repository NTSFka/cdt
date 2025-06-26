package internal

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutableToolNotAvailable(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", nil)

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())
	assert.EqualError(t, tool.Run(Project{}, []string{}), tool.NotFoundError().Error())
}

func TestExecutableTool(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", &Executable{Path: "/bin/tool"})

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}
}

type testTool struct {
	ExecutableTool
}

func TestToolsActiveEmpty(t *testing.T) {
	tools := Tools{}

	assert.Empty(t, tools.Active())
}

func TestToolsActiveNoActive(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("toolId", "", "", nil),
		},
	}

	assert.Empty(t, tools.Active())
}

func TestToolsActive(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("id1", "", "", nil),
		},
		&testTool{
			MakeExecutableTool("id2", "", "", &Executable{Path: "/bin/tool"}),
		},
	}

	active := tools.Active()
	assert.NotEmpty(t, active)
	if assert.Len(t, active, 1) {
		assert.Equal(t, "id2", active[0].Id())
	}
}

func TestToolsGetToolNotFound(t *testing.T) {
	tools := Tools{}

	assert.PanicsWithValue(t, "Tool not found", func() {
		GetTool[*testTool](tools)
	})
}

func TestToolsGetTool(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("toolId", "", "", nil),
		},
	}

	tool := GetTool[*testTool](tools)

	assert.NotNil(t, tool)
	assert.Equal(t, "toolId", tool.Id())
}

func TestPrintToolsEmpty(t *testing.T) {
	tools := Tools{}

	var output bytes.Buffer
	PrintToolList(&output, tools)

	assert.Empty(t, output.String())
}

func TestPrintTools(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("id1", "Tool 1", "", nil),
		},
		&testTool{
			MakeExecutableTool("id2", "Tool 2", "", &Executable{Path: "/bin/tool"}),
		},
	}

	var output bytes.Buffer
	PrintToolList(&output, tools)

	assert.NotEmpty(t, output.String())
}
