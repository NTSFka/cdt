package internal

import (
	"fmt"
	"os"
	"os/exec"
)

// Executable is a structure that stores information about executable
type Executable struct {
	// Command name
	Path string
}

// FindExecutable finds the executable on the system by name
func FindExecutable(name string) *Executable {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil
	}

	return &Executable{Path: path}
}

// Run starts the executable with the given arguments
func (t *Executable) Run(args []string) error {
	command := exec.Command(t.Path, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	fmt.Printf("RUN: %s\n", command.Args)
	return command.Run()
}
