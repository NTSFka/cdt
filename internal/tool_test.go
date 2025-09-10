package internal

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTool_MakeExecutableTool_NotAvailable(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", Tags{}, func() *Executable { return nil })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())
	assert.EqualError(t, tool.Run(context.Background(), RunOptions{}, []string{}), tool.NotFoundError().Error())
}

func TestTool_MakeExecutableTool(t *testing.T) {
	tool := MakeExecutableTool("toolId", "Tool Name", "Tool Info", Tags{}, func() *Executable { return &Executable{Path: "/bin/tool"} })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}
}

func TestTool_NewExecutableTool_NotAvailable(t *testing.T) {
	tool := NewExecutableTool("toolId", "Tool Name", "Tool Info", Tags{}, func() *Executable { return nil })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())
	assert.EqualError(t, tool.Run(context.Background(), RunOptions{}, []string{}), tool.NotFoundError().Error())
}

func TestTool_NewExecutableTool(t *testing.T) {
	tool := NewExecutableTool("toolId", "Tool Name", "Tool Info", Tags{}, func() *Executable { return &Executable{Path: "/bin/tool"} })

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}
}

func TestTool_ExecutableTool_Run(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)
	runtime.On("Id").Return("test")

	tool := MakeExecutableTool("id", "", "", Tags{}, func() *Executable {
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
	runtime.Test(t)
	runtime.On("Id").Return("test")

	tool := MakeExecutableTool("id", "", "", Tags{}, func() *Executable {
		return &Executable{Path: "echo", Runtime: runtime}
	})

	runtime.On("RunExecutable", mock.Anything, RunOptions{
		Directory: "",
		Input:     os.Stdin,
		Output:    os.Stdout,
		Error:     os.Stderr,
	}, "echo", []string{"arg1", "arg2"}).Return(nil)

	err := tool.RunForProject(context.Background(), ProjectInfo{}, []string{"arg1", "arg2"})
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
		NewExecutableTool("toolId", "", "", Tags{}, func() *Executable { return nil }),
	}

	assert.Empty(t, tools.OnlyAvailable())
}

func TestTool_Tools_OnlyAvailable(t *testing.T) {
	tools := Tools{
		NewExecutableTool("id1", "", "", Tags{}, func() *Executable { return nil }),
		NewExecutableTool("id2", "", "", Tags{}, func() *Executable { return &Executable{Path: "/bin/tool"} }),
	}

	active := tools.OnlyAvailable()
	assert.NotEmpty(t, active)
	if assert.Len(t, active, 1) {
		assert.Equal(t, "id2", active[0].Id())
	}
}

func TestTool_Tools_FilterByTags_NotFound(t *testing.T) {
	tools := Tools{
		NewExecutableTool("id1", "", "", Tags{"tag1"}, func() *Executable { return nil }),
		NewExecutableTool("id2", "", "", Tags{"tag2"}, func() *Executable { return &Executable{Path: "/bin/tool"} }),
	}

	active := tools.FilterByTags([]string{"tag0"})
	assert.Empty(t, active)
}

func TestTool_Tools_FilterByTags_Found(t *testing.T) {
	tools := Tools{
		NewExecutableTool("id1", "", "", Tags{"tag1", "tag11"}, func() *Executable { return nil }),
		NewExecutableTool("id2", "", "", Tags{"tag2", "tag21"}, func() *Executable { return &Executable{Path: "/bin/tool"} }),
	}

	data := []struct {
		tags []string
		id   string
	}{
		{
			tags: []string{"tag1"},
			id:   "id1",
		},
		{
			tags: []string{"tag2"},
			id:   "id2",
		},
	}

	for _, d := range data {
		t.Run(d.id, func(t *testing.T) {
			active := tools.FilterByTags(d.tags)
			assert.NotEmpty(t, active)
			if assert.Len(t, active, 1) {
				assert.Equal(t, d.id, active[0].Id())
			}
		})
	}
}

func TestTool_Tools_FilterByTags_Found_MultipleTags(t *testing.T) {
	tools := Tools{
		NewExecutableTool("id1", "", "", Tags{"tag1", "tag2"}, func() *Executable { return nil }),
		NewExecutableTool("id2", "", "", Tags{"tag3", "tag4"}, func() *Executable { return &Executable{Path: "/bin/tool"} }),
	}

	data := []struct {
		tags []string
		id   string
	}{
		{
			tags: []string{"tag1", "tag2"},
			id:   "id1",
		},
		{
			tags: []string{"tag3", "tag4"},
			id:   "id2",
		},
	}

	for _, d := range data {
		t.Run(d.id, func(t *testing.T) {
			active := tools.FilterByTags(d.tags)
			assert.NotEmpty(t, active)
			if assert.Len(t, active, 1) {
				assert.Equal(t, d.id, active[0].Id())
			}
		})
	}
}

func TestTool_Tools_Get_NotFound(t *testing.T) {
	tools := Tools{}

	tool := tools.Get("toolId")

	assert.Nil(t, tool)
}

func TestTool_Tools_Get(t *testing.T) {
	tools := Tools{
		NewExecutableTool("toolId", "", "", Tags{}, nil),
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
		&testTool{MakeExecutableTool("toolId", "", "", Tags{}, nil)},
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

type testRuntime struct {
}

func (t *testRuntime) Id() string {
	return "test"
}

func (t *testRuntime) RunExecutable(_ context.Context, _ RunOptions, _ string, _ []string) error {
	return nil
}

func TestTool_PrintTable(t *testing.T) {
	runtime := testRuntime{}

	tools := Tools{
		NewExecutableTool("id1", "Tool 1", "", Tags{}, func() *Executable { return nil }),
		NewExecutableTool("id2", "Tool 2", "", Tags{}, func() *Executable { return &Executable{Path: "/bin/tool", Runtime: &runtime} }),
	}

	var output bytes.Buffer
	tools.PrintTable(&output)

	assert.NotEmpty(t, output.String())
}
