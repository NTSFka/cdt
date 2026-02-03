package command_test

import (
	"cdt/internal"
	"cdt/internal/command"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtilsParseOptionOutput_Empty(t *testing.T) {
	options := command.ParseOptionOutput[string]("", "def")

	assert.Equal(t, "def", options.Format)
	assert.Nil(t, options.Filename)
}

func TestUtilsParseOptionOutput_OnlyFormat(t *testing.T) {
	options := command.ParseOptionOutput[string]("format1", "def")

	assert.Equal(t, "format1", options.Format)
	assert.Nil(t, options.Filename)
}

func TestUtilsParseOptionOutput_Both(t *testing.T) {
	options := command.ParseOptionOutput[string]("format1:filename2", "def")

	assert.Equal(t, "format1", options.Format)
	require.NotNil(t, options.Filename)
	assert.Equal(t, internal.StrPtr("filename2"), options.Filename)
}
