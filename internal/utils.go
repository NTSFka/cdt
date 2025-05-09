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
