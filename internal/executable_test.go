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

func TestExecutable_Run_Args(t *testing.T) {
	executable := Executable{
		Path: "echo",
		RunFunc: func(ctx context.Context, options RunOptions, path string, args []string) error {
			assert.Equal(t, "echo", path)
			assert.Equal(t, []string{"arg1", "arg2"}, args)

			return nil
		},
	}

	err := executable.Run(context.Background(), RunOptions{}, []string{"arg1", "arg2"})
	assert.NoError(t, err)
}

func TestExecutable_Run_ArgsExtra(t *testing.T) {
	executable := Executable{
		Path: "print",
		Args: []string{"arg1", "arg2"},
		RunFunc: func(ctx context.Context, options RunOptions, path string, args []string) error {
			assert.Equal(t, "print", path)
			assert.Equal(t, []string{"arg1", "arg2", "arg3", "arg4"}, args)

			return nil
		},
	}

	err := executable.Run(context.Background(), RunOptions{}, []string{"arg3", "arg4"})
	assert.NoError(t, err)
}
