package cli

import (
	"cdt/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func runTool(tools internal.Tools, args ...string) error {
	var runArgs []string
	runArgs = append(runArgs, "tool")
	runArgs = append(runArgs, args...)

	return runMainWithTools(tools, runArgs...)
}

func TestToolListEmpty(t *testing.T) {
	err := runTool(internal.Tools{}, "list")

	assert.NoError(t, err)
}

func TestToolListAll(t *testing.T) {
	err := runTool(internal.Tools{}, "list", "--all")

	assert.NoError(t, err)
}

// TODO: check output with some tools
