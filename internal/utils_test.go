package internal

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestUtils_PathExists_NotFound(t *testing.T) {
	assert.False(t, PathExists(filepath.Join(t.TempDir(), "not-found")))
}

func TestUtils_PathExists_Found(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")

	_, err := os.Create(path)
	assert.NoError(t, err)

	assert.True(t, PathExists(path))
}
