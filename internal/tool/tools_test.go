package tool

import (
	"cdt/internal/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitTools(t *testing.T) {
	env := test.NewEnvironment(t)

	tools := InitTools(t.Context(), env)

	assert.NotEmpty(t, tools)
}

func TestInitEnvironmentProviders(t *testing.T) {
	env := test.NewEnvironment(t)

	providers := InitEnvironmentProviders(t.Context(), env)

	assert.NotEmpty(t, providers)
}
