package tool

import (
	"cdt/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTools_SupportedTools_ToTools_Empty(t *testing.T) {
	supportedTools := SupportedTools{}

	tools := supportedTools.ToTools()
	assert.Empty(t, tools)
}

func TestTools_SupportedTools_ToTools(t *testing.T) {
	supportedTools := SupportedTools{
		Go: NewGo(func() *internal.Executable {
			return nil
		}),
	}

	tools := supportedTools.ToTools()
	assert.Len(t, tools, 1)
	assert.Equal(t, supportedTools.Go, tools[0])
}
