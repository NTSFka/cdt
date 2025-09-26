package tool_test

import (
	"testing"

	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestInitTools(t *testing.T) {
	env := test.NewEnvironment(t)

	tools := tool.InitTools(t.Context(), env)

	assert.NotEmpty(t, tools)
}

func TestInitEnvironmentProviders(t *testing.T) {
	env := test.NewEnvironment(t)

	providers := tool.InitEnvironmentProviders(t.Context(), env)

	assert.NotEmpty(t, providers)
}
