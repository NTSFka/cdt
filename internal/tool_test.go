package internal_test

import (
	"bytes"
	"cdt/internal/test"
	"context"
	"os"
	"slices"
	"testing"

	"cdt/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTool_MakeExecutableTool_NotAvailable(t *testing.T) {
	tool := internal.MakeExecutableTool(
		"toolId",
		"Tool Name",
		"Tool Info",
		internal.Tags{},
		test.LazyExecutableNil,
	)

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())
	require.EqualError(
		t,
		tool.Run(t.Context(), internal.RunOptions{}, []string{}),
		tool.NotFoundError().Error(),
	)
}

func TestTool_MakeExecutableTool(t *testing.T) {
	tool := internal.MakeExecutableTool(
		"toolId",
		"Tool Name",
		"Tool Info",
		internal.Tags{},
		test.LazyExecutable("/bin/tool"),
	)

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/tool", executable.Path)
	}
}

func TestTool_NewExecutableTool_NotAvailable(t *testing.T) {
	tool := internal.NewExecutableTool(
		"toolId",
		"Tool Name",
		"Tool Info",
		internal.Tags{},
		test.LazyExecutableNil,
	)

	assert.Equal(t, "toolId", tool.Id())
	assert.Equal(t, "Tool Name", tool.Name())
	assert.Equal(t, "Tool Info", tool.Info())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())
	require.EqualError(
		t,
		tool.Run(t.Context(), internal.RunOptions{}, []string{}),
		tool.NotFoundError().Error(),
	)
}

func TestTool_NewExecutableTool(t *testing.T) {
	tool := internal.NewExecutableTool(
		"toolId",
		"Tool Name",
		"Tool Info",
		internal.Tags{},
		test.LazyExecutable("/bin/tool"),
	)

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

	tool := internal.MakeExecutableTool(
		"id",
		"",
		"",
		internal.Tags{},
		func() (*internal.Executable, error) {
			return &internal.Executable{Path: "echo", Runtime: runtime}, nil
		},
	)

	runtime.On("RunExecutable", mock.Anything, internal.RunOptions{}, "echo", []string{"arg1", "arg2"}).
		Return(nil)

	err := tool.Run(t.Context(), internal.RunOptions{}, []string{"arg1", "arg2"})
	require.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestTool_ExecutableTool_RunForProject(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)
	runtime.On("Id").Return("test")

	tool := internal.MakeExecutableTool(
		"id",
		"",
		"",
		internal.Tags{},
		func() (*internal.Executable, error) {
			return &internal.Executable{Path: "echo", Runtime: runtime}, nil
		},
	)

	runtime.On("RunExecutable", mock.Anything, internal.RunOptions{
		Directory: "",
		Input:     os.Stdin,
		Output:    os.Stdout,
		Error:     os.Stderr,
	}, "echo", []string{"arg1", "arg2"}).Return(nil)

	err := tool.RunForProject(t.Context(), internal.ProjectInfo{}, []string{"arg1", "arg2"})
	require.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestTool_ExecutableTool_RunForProjectWithOutput(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)
	runtime.On("Id").Return("test")

	tool := internal.MakeExecutableTool(
		"id",
		"",
		"",
		internal.Tags{},
		func() (*internal.Executable, error) {
			return &internal.Executable{Path: "echo", Runtime: runtime}, nil
		},
	)

	outputFilename := t.TempDir() + "/output.txt"

	runtime.On("RunExecutable", mock.Anything, mock.MatchedBy(func(options internal.RunOptions) bool {
		require.IsType(t, &os.File{}, options.Output)
		assert.Equal(t, outputFilename, options.Output.(*os.File).Name())

		return true
	}), "echo", []string{"arg1", "arg2"}).
		Return(nil)

	err := tool.RunForProjectWithOutput(
		t.Context(),
		internal.ProjectInfo{},
		outputFilename,
		[]string{"arg1", "arg2"},
	)
	require.NoError(t, err)

	runtime.AssertExpectations(t)
}

func TestTool_ExecutableTool_RunForProjectWithOutput_FileFailed(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)

	tool := internal.MakeExecutableTool(
		"id",
		"",
		"",
		internal.Tags{},
		func() (*internal.Executable, error) {
			return &internal.Executable{Path: "echo", Runtime: runtime}, nil
		},
	)

	outputFilename := t.TempDir() + "/non-existent-directory/output.txt"

	err := tool.RunForProjectWithOutput(
		t.Context(),
		internal.ProjectInfo{},
		outputFilename,
		[]string{"arg1", "arg2"},
	)
	require.ErrorContains(t, err, "failed to create output file")

	runtime.AssertExpectations(t)
}

func TestTool_ExecutableTool_RunForProjectWithEnv(t *testing.T) {
	runtime := &testExecutableRuntime{}
	runtime.Test(t)
	runtime.On("Id").Return("test")

	tool := internal.MakeExecutableTool(
		"id",
		"",
		"",
		internal.Tags{},
		func() (*internal.Executable, error) {
			return &internal.Executable{Path: "echo", Runtime: runtime}, nil
		},
	)

	runtime.On("RunExecutable", mock.Anything, mock.MatchedBy(func(options internal.RunOptions) bool {
		if !slices.Contains(options.Env, "VAR1=value1") {
			return false
		}

		if !slices.Contains(options.Env, "VAR2=value2") {
			return false
		}

		return true
	}), "echo", []string{"arg1", "arg2"}).
		Return(nil)

	err := tool.RunForProjectWithEnv(
		t.Context(),
		internal.ProjectInfo{},
		[]string{"VAR1=value1", "VAR2=value2"},
		[]string{"arg1", "arg2"},
	)
	require.NoError(t, err)

	runtime.AssertExpectations(t)
}

type testTool struct {
	internal.ExecutableTool
}

func TestTool_Tools_OnlyAvailable_Empty(t *testing.T) {
	tools := internal.Tools{}

	assert.Empty(t, tools.OnlyAvailable())
}

func TestTool_Tools_OnlyAvailable_NotAvailable(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("toolId", "", "", internal.Tags{}, test.LazyExecutableNil),
	}

	assert.Empty(t, tools.OnlyAvailable())
}

func TestTool_Tools_OnlyAvailable(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("id1", "", "", internal.Tags{}, test.LazyExecutableNil),
		internal.NewExecutableTool(
			"id2",
			"",
			"",
			internal.Tags{},
			test.LazyExecutable("/bin/tool"),
		),
	}

	active := tools.OnlyAvailable()
	assert.NotEmpty(t, active)
	require.Len(t, active, 1)
	assert.Equal(t, "id2", active[0].Id())
}

func TestTool_Tools_FilterByTags_NotFound(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool(
			"id1",
			"",
			"",
			internal.Tags{"tag1"},
			test.LazyExecutableNil,
		),
		internal.NewExecutableTool(
			"id2",
			"",
			"",
			internal.Tags{"tag2"},
			test.LazyExecutable("/bin/tool"),
		),
	}

	active := tools.FilterByTags([]string{"tag0"})
	assert.Empty(t, active)
}

func TestTool_Tools_FilterByTags_Found(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool(
			"id1",
			"",
			"",
			internal.Tags{"tag1", "tag11"},
			test.LazyExecutableNil,
		),
		internal.NewExecutableTool(
			"id2",
			"",
			"",
			internal.Tags{"tag2", "tag21"},
			test.LazyExecutable("/bin/tool"),
		),
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
			require.Len(t, active, 1)
			assert.Equal(t, d.id, active[0].Id())
		})
	}
}

func TestTool_Tools_FilterByTags_Found_MultipleTags(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool(
			"id1",
			"",
			"",
			internal.Tags{"tag1", "tag2"},
			test.LazyExecutableNil,
		),
		internal.NewExecutableTool(
			"id2",
			"",
			"",
			internal.Tags{"tag3", "tag4"},
			test.LazyExecutable("/bin/tool"),
		),
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
			require.Len(t, active, 1)
			assert.Equal(t, d.id, active[0].Id())
		})
	}
}

func TestTool_Tools_Get_NotFound(t *testing.T) {
	tools := internal.Tools{}

	tool := tools.Get("toolId")

	assert.Nil(t, tool)
}

func TestTool_Tools_Get(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("toolId", "", "", internal.Tags{}, nil),
	}

	tool := tools.Get("toolId")

	assert.NotNil(t, tool)
	assert.Equal(t, "toolId", tool.Id())
}

func TestTool_Tools_GetTool_NotFound(t *testing.T) {
	tools := internal.Tools{}

	assert.PanicsWithValue(t, "Tool not found", func() {
		internal.GetTool[*testTool](tools)
	})
}

func TestTool_Tools_GetTool(t *testing.T) {
	tools := internal.Tools{
		&testTool{internal.MakeExecutableTool("toolId", "", "", internal.Tags{}, nil)},
	}

	tool := internal.GetTool[*testTool](tools)

	assert.NotNil(t, tool)
	assert.Equal(t, "toolId", tool.Id())
}

func TestTool_PrintTable_Empty(t *testing.T) {
	tools := internal.Tools{}

	var output bytes.Buffer

	tools.PrintTable(&output)

	assert.Empty(t, output.String())
}

type testRuntime struct {
}

func (t *testRuntime) Id() string {
	return "test"
}

func (t *testRuntime) RunExecutable(
	_ context.Context,
	_ internal.RunOptions,
	_ string,
	_ []string,
) error {
	return nil
}

func TestTool_PrintTable(t *testing.T) {
	runtime := testRuntime{}

	tools := internal.Tools{
		internal.NewExecutableTool(
			"id1",
			"Tool 1",
			"",
			internal.Tags{},
			test.LazyExecutableNil,
		),
		internal.NewExecutableTool(
			"id2",
			"Tool 2",
			"",
			internal.Tags{},
			func() (*internal.Executable, error) {
				return &internal.Executable{Path: "/bin/tool", Runtime: &runtime}, nil
			},
		),
	}

	var output bytes.Buffer

	tools.PrintTable(&output)

	assert.NotEmpty(t, output.String())
}

func TestTool_GetToolPath_NotFound(t *testing.T) {
	options := internal.DetectOptions{}

	path := options.GetToolPath("tool1", "default-tool")
	assert.Equal(t, "default-tool", path)
}

func TestTool_GetToolPath_Found(t *testing.T) {
	options := internal.DetectOptions{
		ToolsPaths: map[string]string{
			"tool1": "tool1-path",
		},
	}

	path := options.GetToolPath("tool1", "default-tool")
	assert.Equal(t, "tool1-path", path)
}
