package internal

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvironment_SystemEnvironment_FindExecutable_NotFound(t *testing.T) {
	executable := SystemEnvironment.FindExecutable("tool-not-found")

	assert.Nil(t, executable)
}

func TestEnvironment_SystemEnvironment_FindExecutable(t *testing.T) {
	executable := SystemEnvironment.FindExecutable("echo")

	if assert.NotNil(t, executable) {
		assert.NotNil(t, executable.RunFunc)
		assert.Contains(t, executable.Path, "echo")
	}
}

func TestEnvironment_SystemEnvironment_RunExecutable(t *testing.T) {
	buffer := bytes.Buffer{}
	options := RunOptions{
		Directory: ".",
		Output:    &buffer,
		Error:     nil,
	}

	err := SystemEnvironment.RunExecutable(context.Background(), options, "echo", []string{"test"})
	assert.NoError(t, err)
	assert.Equal(t, "test\n", buffer.String())
}
