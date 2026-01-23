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

// TraceHandle capture tracing start of some call and allowing to finish the trace later via Done().
type TraceHandle struct {
	name   string
	logger *slog.Logger
	start  time.Time
}

// TraceStart starts tracing of some call.
func TraceStart(ctx context.Context, name string, args ...any) TraceHandle {
	logger := slog.With(args...)

	logger.DebugContext(ctx, indent()+"+ "+name)

	traceIndent++

	start := time.Now()

	return TraceHandle{
		name:   name,
		logger: logger,
		start:  start,
	}
}

func (t TraceHandle) Done(ctx context.Context, args ...any) {
	traceIndent--

	t.logger.DebugContext(
		ctx,
		indent()+"- "+t.name,
		append(args, "duration", time.Since(t.start))...)
}

// Trace captures a function call by storing start and end.
func Trace[R any](ctx context.Context, name string, function func() R, args ...any) R {
	trace := TraceStart(ctx, name, args...)

	res := function()

	trace.Done(ctx, "result", res)

	return res
}

// TraceErr captures a function call by storing start and end.
func TraceErr[R any](
	ctx context.Context,
	name string,
	function func() (R, error),
	args ...any,
) (R, error) {
	trace := TraceStart(ctx, name, args...)

	res, err := function()

	trace.Done(ctx, "result", res, "error", err)

	return res, err
}
