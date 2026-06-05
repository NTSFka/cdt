package internal_test

import (
	"bytes"
	"cdt/internal"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLogger(t *testing.T) *bytes.Buffer {
	t.Helper()

	buf := new(bytes.Buffer)
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	logger := slog.New(slog.NewTextHandler(buf, opts))

	oldLogger := slog.Default()

	slog.SetDefault(logger)

	t.Cleanup(func() {
		slog.SetDefault(oldLogger)
	})

	return buf
}

func TestLog_TraceStart(t *testing.T) {
	buffer := setupLogger(t)

	trace := internal.TraceStart(t.Context(), "test")

	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.NotContains(t, buffer.String(), "msg=\"- test\"")
	assert.NotContains(t, buffer.String(), "duration=")

	trace.Done(t.Context(), "result", 67)

	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=67")
}

func TestLog_TraceStart_WithArgs(t *testing.T) {
	buffer := setupLogger(t)

	trace := internal.TraceStart(t.Context(), "test", "arg1", 1, "arg2", "Test")

	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.NotContains(t, buffer.String(), "msg=\"- test\"")
	assert.NotContains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "arg1=1")
	assert.Contains(t, buffer.String(), "arg2=Test")

	trace.Done(t.Context(), "result", 67)

	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=67")
}

func TestLog_Trace(t *testing.T) {
	buffer := setupLogger(t)

	res := internal.Trace(t.Context(), "test", func() int { return 42 })

	assert.Equal(t, 42, res)
	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=42")
}

func TestLog_Trace_WithArgs(t *testing.T) {
	buffer := setupLogger(t)

	res := internal.Trace(
		t.Context(),
		"test",
		func() int { return 42 },
		"arg1",
		66,
		"arg2",
		"Hello",
	)

	assert.Equal(t, 42, res)
	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=42")
	assert.Contains(t, buffer.String(), "arg1=66")
	assert.Contains(t, buffer.String(), "arg2=Hello")
}

func TestLog_Trace_Nested(t *testing.T) {
	buffer := setupLogger(t)

	res := internal.Trace(t.Context(), "test", func() int {
		return internal.Trace(t.Context(), "nested", func() int { return 42 })
	})

	assert.Equal(t, 42, res)
	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\" + nested\"")
	assert.Contains(t, buffer.String(), "msg=\" - nested\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=42")
}

func TestLog_TraceErr_NoError(t *testing.T) {
	buffer := setupLogger(t)

	res, err := internal.TraceErr(t.Context(), "test", func() (int, error) { return 42, nil })

	assert.Equal(t, 42, res)
	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=42")
}

func TestLog_TraceErr_WithArgs(t *testing.T) {
	buffer := setupLogger(t)

	res, err := internal.TraceErr(
		t.Context(),
		"test",
		func() (int, error) { return 42, nil },
		"arg1",
		22,
		"arg2",
		"World",
	)

	assert.Equal(t, 42, res)
	require.NoError(t, err)
	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=42")
	assert.Contains(t, buffer.String(), "arg1=22")
	assert.Contains(t, buffer.String(), "arg2=World")
}

func TestLog_TraceErr_Error(t *testing.T) {
	buffer := setupLogger(t)

	res, err := internal.TraceErr(
		t.Context(),
		"test",
		func() (int, error) { return 2, errors.New("failed") },
	)

	assert.Equal(t, 2, res)
	require.EqualError(t, err, "failed")
	assert.Contains(t, buffer.String(), "msg=\"+ test\"")
	assert.Contains(t, buffer.String(), "msg=\"- test\"")
	assert.Contains(t, buffer.String(), "duration=")
	assert.Contains(t, buffer.String(), "result=2")
}
