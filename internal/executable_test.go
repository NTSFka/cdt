package internal

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutable_Run(t *testing.T) {
	executable := Executable{Path: "echo", RunFunc: func(ctx context.Context, options RunOptions, path string, args []string) error {
		return nil
	}}

	err := executable.Run(context.Background(), RunOptions{}, []string{})
	assert.NoError(t, err)
}

func TestExecutable_Run_Failed(t *testing.T) {
	executable := Executable{Path: "echo", RunFunc: func(ctx context.Context, options RunOptions, path string, args []string) error {
		return errors.New("failed")
	}}

	err := executable.Run(context.Background(), RunOptions{}, []string{})
	assert.EqualError(t, err, "failed")
}
