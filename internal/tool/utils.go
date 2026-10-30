package tool

import (
	"errors"
	"os/exec"
	"slices"
)

// appendFiles appends files to the arguments if they are not nil or empty.
func appendFiles(args []string, filenames *[]string, all *string) []string {
	if filenames != nil && len(*filenames) > 0 {
		args = append(args, *filenames...)
	} else if all != nil {
		args = append(args, *all)
	}

	return args
}

// checkExitCode checks if the error is an exec.ExitError and if its exit code is in the acceptance list.
func checkExitCode(err error, accept []int) bool {
	exitErr, ok := errors.AsType[*exec.ExitError](err)

	if !ok {
		return false
	}

	return slices.Contains(accept, exitErr.ExitCode())
}
