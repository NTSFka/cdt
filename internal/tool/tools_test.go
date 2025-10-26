package tool_test

import (
	"cdt/internal"
	"testing"

	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestInitTools(t *testing.T) {
	env := test.NewEnvironment(t)

	tools := tool.InitTools(t.Context(), internal.ConfigTools{}, env)

	assert.NotEmpty(t, tools)
}

func TestInitEnvironmentProviders(t *testing.T) {
	env := test.NewEnvironment(t)

	providers := tool.InitEnvironmentProviders(t.Context(), internal.ConfigTools{}, env)

	assert.NotEmpty(t, providers)
}
