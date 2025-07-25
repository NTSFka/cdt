package internal

import (
	"errors"
	"os"
)

// PathExists checks if a path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// Assert panics when the condition is not true
func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

// StrPtr returns a pointer to a given string
func StrPtr(s string) *string {
	return &s
}
