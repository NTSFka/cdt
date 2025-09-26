package internal

import (
	"errors"
	"io"
	"os"

	"golang.org/x/term"
)

// PathExists checks if a path exists.
func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// Assert panics when the condition is not true.
func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

// StrPtr returns a pointer to a given string.
func StrPtr(s string) *string {
	return &s
}

// DetectTermWidth detects the terminal width from a writer if possible.
func DetectTermWidth(writer io.Writer) *int {
	if f, ok := writer.(interface{ Fd() uintptr }); ok && term.IsTerminal(int(f.Fd())) {
		if w, _, err := term.GetSize(int(f.Fd())); err == nil && w > 0 {
			return &w
		}
	} else if term.IsTerminal(int(os.Stdout.Fd())) {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
			return &w
		}
	}

	return nil
}
