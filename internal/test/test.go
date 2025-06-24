package test

import (
	"io"
	"os"
)

// CaptureOutput captures stdout and returns it as string
func CaptureOutput(f func() error) (string, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig

	if err := w.Close(); err != nil {
		return "", err
	}

	out, _ := io.ReadAll(r)
	return string(out), err
}
