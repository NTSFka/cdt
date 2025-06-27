package internal

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExecutable_NewRunContext(t *testing.T) {
	ctx := NewRunContext("project")

	assert.Equal(t, "project", ctx.Directory)
	assert.Equal(t, os.Stdout, ctx.Output)
	assert.Equal(t, os.Stderr, ctx.Error)
}

func TestExecutable_Run(t *testing.T) {
	executable := Executable{Path: "echo", RunFunc: func(ctx RunContext, path string, args []string) error {
		return nil
	}}

	err := executable.Run(RunContext{}, []string{})
	assert.NoError(t, err)
}

func TestExecutable_Run_Failed(t *testing.T) {
	executable := Executable{Path: "echo", RunFunc: func(ctx RunContext, path string, args []string) error {
		return errors.New("failed")
	}}

	err := executable.Run(RunContext{}, []string{})
	assert.EqualError(t, err, "failed")
}
