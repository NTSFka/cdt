package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTool_List_Empty(t *testing.T) {
	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{}, "list")

	require.NoError(t, err)
}

func TestTool_List_Tags(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", internal.Tags{"tag1"}, func() *internal.Executable {
			return nil
		}),
	}

	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{Tools: tools}, "list", "--tag", "tag1")

	require.NoError(t, err)
}

func TestTool_ListAll_Empty(t *testing.T) {
	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{}, "list", "--all")

	require.NoError(t, err)
}

func TestTool_ListAll(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", internal.Tags{}, func() *internal.Executable {
			return nil
		}),
	}

	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{Tools: tools}, "list", "--all")

	require.NoError(t, err)
}

func TestTool_Run_Unknown(t *testing.T) {
	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{}, "run", "tool")

	require.EqualError(t, err, "tool 'tool' not found")
}

func TestTool_Run_Unavailable(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", internal.Tags{}, func() *internal.Executable {
			return nil
		}),
	}

	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{Tools: tools}, "run", "tool")

	require.EqualError(t, err, "tool 'tool' failed: Tool is not installed on the system")
}

func TestTool_Run_Success(t *testing.T) {
	exec := test.NewExecutable(t)

	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", internal.Tags{}, exec.LazyExecutable("/usr/bin/tool")),
	}

	exec.OnRun("/usr/bin/tool", []string{}).
		Return(nil)

	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{Tools: tools}, "run", "tool")

	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestTool_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", internal.Tags{}, exec.LazyExecutable("/usr/bin/tool")),
	}

	exec.OnRun("/usr/bin/tool", []string{}).
		Return(errors.New("failed"))

	err := test.RunCommand(t.Context(), NewToolCommand(), internal.Context{Tools: tools}, "run", "tool")

	require.EqualError(t, err, "tool 'tool' failed: failed")

	exec.AssertExpectations(t)
}
