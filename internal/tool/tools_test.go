package tool

import (
	"cdt/internal/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitTools(t *testing.T) {
	env := test.NewEnvironment(t)

	tools := InitTools(env)

	assert.NotEmpty(t, tools)
}

func TestInitEnvironmentProviders(t *testing.T) {
	t.Skip("FIXME")

	env := test.NewEnvironment(t)

	providers := InitEnvironmentProviders(env)

	assert.NotEmpty(t, providers)
}
