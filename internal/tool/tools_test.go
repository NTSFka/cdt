package tool_test

import (
	"testing"

	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestInitTools(t *testing.T) {
	env := test.NewEnvironment(t)

	tools := tool.InitTools(t.Context(), tool.DetectOptions{Environment: env})

	assert.NotEmpty(t, tools)
}

func TestInitEnvironmentProviders(t *testing.T) {
	env := test.NewEnvironment(t)

	providers := tool.InitEnvironmentProviders(t.Context(), tool.DetectOptions{Environment: env})

	assert.NotEmpty(t, providers)
}

func TestDetectOptions_GetToolPath_NotFound(t *testing.T) {
	options := tool.DetectOptions{}

	path := options.GetToolPath("tool1", "default-tool")
	assert.Equal(t, "default-tool", path)
}

func TestDetectOptions_GetToolPath_Found(t *testing.T) {
	options := tool.DetectOptions{
		ToolsPaths: map[string]string{
			"tool1": "tool1-path",
		},
	}

	path := options.GetToolPath("tool1", "default-tool")
	assert.Equal(t, "tool1-path", path)
}
