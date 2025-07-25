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

func TestAssert_Success(t *testing.T) {
	assert.NotPanics(t, func() {
		Assert(true, "success")
	})
}

func TestAssert_Fails(t *testing.T) {
	assert.Panics(t, func() {
		Assert(false, "failed")
	}, "failed")
}

func TestStrPtr(t *testing.T) {
	s := "test"

	assert.Equal(t, &s, StrPtr(s))
}
