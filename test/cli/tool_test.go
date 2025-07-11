package cli

import (
	"cdt/internal"
	"cdt/internal/test"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func runTool(tools internal.Tools, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "tool")
	runArgs = append(runArgs, args...)

	return runMainWithTools(tools, runArgs...)
}

func TestTool_List_Empty(t *testing.T) {
	err := runTool(internal.Tools{}, "list")

	assert.NoError(t, err)
}

func TestTool_ListAll_Empty(t *testing.T) {
	err := runTool(internal.Tools{}, "list", "--all")

	assert.NoError(t, err)
}

func TestTool_ListAll(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", func() *internal.Executable {
			return nil
		}),
	}

	err := runTool(tools, "list", "--all")

	assert.NoError(t, err)
}

func TestTool_Run_Unknown(t *testing.T) {
	err := runTool(internal.Tools{}, "run", "tool")

	assert.EqualError(t, err, "tool 'tool' not found")
}

func TestTool_Run_Unavailable(t *testing.T) {
	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", func() *internal.Executable {
			return nil
		}),
	}

	err := runTool(tools, "run", "tool")

	assert.EqualError(t, err, "tool 'tool' failed: Tool is not installed on the system")
}

func TestTool_Run_Success(t *testing.T) {
	runMock := test.Executable{}

	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", runMock.LazyExecutable("/usr/bin/tool")),
	}

	runMock.OnRun("/usr/bin/tool", []string{}).
		Return(nil)

	err := runTool(tools, "run", "tool")

	assert.NoError(t, err)
}

func TestTool_Run_Failed(t *testing.T) {
	runMock := test.Executable{}

	tools := internal.Tools{
		internal.NewExecutableTool("tool", "Tool", "Some tool", runMock.LazyExecutable("/usr/bin/tool")),
	}

	runMock.OnRun("/usr/bin/tool", []string{}).
		Return(errors.New("failed"))

	err := runTool(tools, "run", "tool")

	assert.EqualError(t, err, "tool 'tool' failed: failed")
}
