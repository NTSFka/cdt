package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtils_PathExists_NotFound(t *testing.T) {
	assert.False(t, PathExists(filepath.Join(t.TempDir(), "not-found")))
}

func TestUtils_PathExists_Found(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file")

	_, err := os.Create(path)
	require.NoError(t, err)

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

func TestDetectTermWidth_NotTerminal(t *testing.T) {
	writer := &bytes.Buffer{}

	width := DetectTermWidth(writer)

	assert.Nil(t, width)
}

func TestDetectTermWidth_StdOut(t *testing.T) {
	width := DetectTermWidth(os.Stdout)

	assert.Nil(t, width)
}
