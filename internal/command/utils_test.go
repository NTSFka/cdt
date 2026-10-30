package command_test

import (
	"cdt/internal/command"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtilsParseOptionOutput_Empty(t *testing.T) {
	options := command.ParseOptionOutput("")

	assert.Empty(t, options.Format)
	assert.Nil(t, options.Filename)
}

func TestUtilsParseOptionOutput_OnlyFormat(t *testing.T) {
	options := command.ParseOptionOutput("format1")

	assert.Equal(t, "format1", options.Format)
	assert.Nil(t, options.Filename)
}

func TestUtilsParseOptionOutput_Both(t *testing.T) {
	options := command.ParseOptionOutput("format1:filename2")

	assert.Equal(t, "format1", options.Format)
	require.NotNil(t, options.Filename)
	assert.Equal(t, new("filename2"), options.Filename)
}
