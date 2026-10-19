package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CreateFile creates a file in the given directory with the given name.
func CreateFile(t *testing.T, directory string, name string) {
	file, err := os.Create(filepath.Join(directory, name))

	require.NoError(t, err)
	assert.NoError(t, file.Close())
}

// CreateDirectory creates a directory in the given directory with the given name.
func CreateDirectory(t *testing.T, directory string, name string) {
	err := os.Mkdir(filepath.Join(directory, name), 0755)

	require.NoError(t, err)
}
