package internal

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTool_ExecutableTool_NotAvailable(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", func() *Executable { return nil })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())
	assert.EqualError(t, tool.Run(Project{}, []string{}), tool.NotFoundError().Error())
}

func TestTool_ExecutableTool(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", func() *Executable { return &Executable{Path: "/bin/tool"} })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}
}

func TestTool_ExecutableTool_Run(t *testing.T) {
	var runMock mock.Mock

	tool := MakeExecutableTool("id", "", "", func() *Executable {
		return &Executable{Path: "echo", RunFunc: func(ctx context.Context, options RunOptions, path string, args []string) error {
			return runMock.Called(ctx, path, args).Error(0)
		}}
	})

	runMock.On("1", mock.Anything, "echo", []string{"arg1", "arg2"}).Return(nil)

	err := tool.Run(Project{}, []string{"arg1", "arg2"})
	assert.NoError(t, err)
}

type testTool struct {
	ExecutableTool
}

func TestTool_Tools_Active_Empty(t *testing.T) {
	tools := Tools{}

	assert.Empty(t, tools.Active())
}

func TestTool_Tools_Active_NoActive(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("toolId", "", "", func() *Executable { return nil }),
		},
	}

	assert.Empty(t, tools.Active())
}

func TestTool_Tools_Active(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("id1", "", "", func() *Executable { return nil }),
		},
		&testTool{
			MakeExecutableTool("id2", "", "", func() *Executable { return &Executable{Path: "/bin/tool"} }),
		},
	}

	active := tools.Active()
	assert.NotEmpty(t, active)
	if assert.Len(t, active, 1) {
		assert.Equal(t, "id2", active[0].Id())
	}
}

func TestTool_Tools_GetTool_NotFound(t *testing.T) {
	tools := Tools{}

	assert.PanicsWithValue(t, "Tool not found", func() {
		GetTool[*testTool](tools)
	})
}

func TestTool_Tools_GetTool(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("toolId", "", "", nil),
		},
	}

	tool := GetTool[*testTool](tools)

	assert.NotNil(t, tool)
	assert.Equal(t, "toolId", tool.Id())
}

func TestTool_PrintTools_Empty(t *testing.T) {
	tools := Tools{}

	var output bytes.Buffer
	PrintToolList(&output, tools)

	assert.Empty(t, output.String())
}

func TestTool_PrintTools(t *testing.T) {
	tools := Tools{
		&testTool{
			MakeExecutableTool("id1", "Tool 1", "", func() *Executable { return nil }),
		},
		&testTool{
			MakeExecutableTool("id2", "Tool 2", "", func() *Executable { return &Executable{Path: "/bin/tool"} }),
		},
	}

	var output bytes.Buffer
	PrintToolList(&output, tools)

	assert.NotEmpty(t, output.String())
}
