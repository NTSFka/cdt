package command

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTool_List_Empty(t *testing.T) {
	err := test.RunCommand(NewToolCommand(), internal.Context{}, "list")

	assert.NoError(t, err)
}

func TestTool_ListAll_Empty(t *testing.T) {
	err := test.RunCommand(NewToolCommand(), internal.Context{}, "list", "--all")

	assert.NoError(t, err)
}

func TestTool_ListAll(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", func() *internal.Executable {
			return nil
		}),
	}

	err := test.RunCommand(NewToolCommand(), internal.Context{Tools: tools}, "list", "--all")

	assert.NoError(t, err)
}

func TestTool_Run_Unknown(t *testing.T) {
	err := test.RunCommand(NewToolCommand(), internal.Context{}, "run", "tool")

	assert.EqualError(t, err, "tool 'tool' not found")
}

func TestTool_Run_Unavailable(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", func() *internal.Executable {
			return nil
		}),
	}

	err := test.RunCommand(NewToolCommand(), internal.Context{Tools: tools}, "run", "tool")

	assert.EqualError(t, err, "tool 'tool' failed: Tool is not installed on the system")
}

func TestTool_Run_Success(t *testing.T) {
	exec := test.NewExecutable(t)

	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", exec.LazyExecutable("/usr/bin/tool")),
	}

	exec.OnRun("/usr/bin/tool", []string{}).
		Return(nil)

	err := test.RunCommand(NewToolCommand(), internal.Context{Tools: tools}, "run", "tool")

	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestTool_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", exec.LazyExecutable("/usr/bin/tool")),
	}

	exec.OnRun("/usr/bin/tool", []string{}).
		Return(errors.New("failed"))

	err := test.RunCommand(NewToolCommand(), internal.Context{Tools: tools}, "run", "tool")

	assert.EqualError(t, err, "tool 'tool' failed: failed")

	exec.AssertExpectations(t)
}
