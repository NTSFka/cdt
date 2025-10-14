package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Infof prints CDT action.
func Infof(format string, a ...any) {
	_, _ = color.New(color.FgCyan).Printf("[cdt] %v\n", fmt.Sprintf(format, a...))
}

// Debug logs debug message.
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

var traceIndent = 0

func indent() string {
	return strings.Repeat(" ", traceIndent)
}

// Trace captures a function call by storing start and end.
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

// TraceErr captures a function call by storing start and end.
func TraceErr[R any](
	ctx context.Context,
	name string,
	function func() (R, error),
	args ...any,
) (R, error) {
	logger := slog.With(args...)

	logger.DebugContext(ctx, indent()+"+ "+name)

	traceIndent++

	start := time.Now()

	res, err := function()

	traceIndent--

	logger.DebugContext(
		ctx,
		indent()+"- "+name,
		"result",
		res,
		"error",
		err,
		"duration",
		time.Since(start),
	)

	return res, err
}
