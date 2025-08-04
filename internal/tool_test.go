package internal

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestTool_MakeExecutableTool_NotAvailable(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", func() *Executable { return nil })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())
	assert.EqualError(t, tool.Run(context.Background(), RunOptions{}, []string{}), tool.NotFoundError().Error())
}

func TestTool_MakeExecutableTool(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", func() *Executable { return &Executable{Path: "/bin/tool"} })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}
}

func TestTool_NewExecutableTool_NotAvailable(t *testing.T) {
	tool := NewExecutableTool("toolId", "Tool Name", "Tool Info", func() *Executable { return nil })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.ExecutablePath())
	assert.EqualError(t, tool.Run(context.Background(), RunOptions{}, []string{}), tool.NotFoundError().Error())
}

func TestTool_NewExecutableTool(t *testing.T) {
	tool := NewExecutableTool("toolId", "Tool Name", "Tool Info", func() *Executable { return &Executable{Path: "/bin/tool"} })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())

	if path := tool.ExecutablePath(); assert.NotNil(t, path) {
		assert.Equal(t, "/bin/tool", *path)
	}
}

func TestTool_ExecutableTool_Run(t *testing.T) {
	runtime := &testExecutableRuntime{}

	tool := MakeExecutableTool("id", "", "", func() *Executable {
		return &Executable{Path: "echo", Runtime: runtime}
	})

	runtime.On("RunExecutable", mock.Anything, RunOptions{}, "echo", []string{"arg1", "arg2"}).
		Return(nil)

	err := tool.Run(context.Background(), RunOptions{}, []string{"arg1", "arg2"})
	assert.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestTool_ExecutableTool_RunForProject(t *testing.T) {
	runtime := &testExecutableRuntime{}

	tool := MakeExecutableTool("id", "", "", func() *Executable {
		return &Executable{Path: "echo", Runtime: runtime}
	})

	runtime.On("RunExecutable", mock.Anything, RunOptions{
		Directory: "",
		Input:     os.Stdin,
		Output:    os.Stdout,
		Error:     os.Stderr,
	}, "echo", []string{"arg1", "arg2"}).Return(nil)

	err := tool.RunForProject(Project{}, []string{"arg1", "arg2"})
	assert.NoError(t, err)

	runtime.AssertExpectations(t)
}

type testTool struct {
	ExecutableTool
}

func TestTool_Tools_OnlyAvailable_Empty(t *testing.T) {
	tools := Tools{}

	assert.Empty(t, tools.OnlyAvailable())
}

func TestTool_Tools_OnlyAvailable_NotAvailable(t *testing.T) {
	tools := Tools{
		NewExecutableTool("toolId", "", "", func() *Executable { return nil }),
	}

	assert.Empty(t, tools.OnlyAvailable())
}

func TestTool_Tools_OnlyAvailable(t *testing.T) {
	tools := Tools{
		NewExecutableTool("id1", "", "", func() *Executable { return nil }),
		NewExecutableTool("id2", "", "", func() *Executable { return &Executable{Path: "/bin/tool"} }),
	}

	active := tools.OnlyAvailable()
	assert.NotEmpty(t, active)
	if assert.Len(t, active, 1) {
		assert.Equal(t, "id2", active[0].Id())
	}
}

func TestTool_Tools_Get_NotFound(t *testing.T) {
	tools := Tools{}

	tool := tools.Get("toolId")

	assert.Nil(t, tool)
}

func TestTool_Tools_Get(t *testing.T) {
	tools := Tools{
		NewExecutableTool("toolId", "", "", nil),
	}

	tool := tools.Get("toolId")

	assert.NotNil(t, tool)
	assert.Equal(t, "toolId", tool.Id())
}

func TestTool_Tools_GetTool_NotFound(t *testing.T) {
	tools := Tools{}

	assert.PanicsWithValue(t, "Tool not found", func() {
		GetTool[*testTool](tools)
	})
}

func TestTool_Tools_GetTool(t *testing.T) {
	tools := Tools{
		&testTool{MakeExecutableTool("toolId", "", "", nil)},
	}

	tool := GetTool[*testTool](tools)

	assert.NotNil(t, tool)
	assert.Equal(t, "toolId", tool.Id())
}

func TestTool_PrintTable_Empty(t *testing.T) {
	tools := Tools{}

	var output bytes.Buffer
	tools.PrintTable(&output)

	assert.Empty(t, output.String())
}

func TestTool_PrintTable(t *testing.T) {
	tools := Tools{
		NewExecutableTool("id1", "Tool 1", "", func() *Executable { return nil }),
		NewExecutableTool("id2", "Tool 2", "", func() *Executable { return &Executable{Path: "/bin/tool"} }),
	}

	var output bytes.Buffer
	tools.PrintTable(&output)

	assert.NotEmpty(t, output.String())
}
