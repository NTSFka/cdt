package internal

import (
	"context"
	"log/slog"
	"strings"
	"time"
)

var traceIndent = 0

func indent() string {
	return strings.Repeat(" ", traceIndent)
}

// Trace captures a function call by storing start and end
func Trace[R any](ctx context.Context, name string, function func() R, args ...any) R {
	logger := slog.With(args...)

	logger.DebugContext(ctx, indent()+"+ "+name)

	traceIndent++

	start := time.Now()

	res := function()

	traceIndent--

	logger.DebugContext(ctx, indent()+"- "+name, "result", res, "duration", time.Since(start))

	return res
}
