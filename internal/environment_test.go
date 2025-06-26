package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSystemEnvironmentFindExecutableNotFound(t *testing.T) {
	executable := SystemEnvironment.FindExecutable("tool-not-found")

	assert.Nil(t, executable)
}

func TestSystemEnvironmentFindExecutable(t *testing.T) {
	executable := SystemEnvironment.FindExecutable("echo")

	if assert.NotNil(t, executable) {
		assert.NotNil(t, executable.RunFunc)
		assert.Contains(t, executable.Path, "echo")
	}
}
