package internal

import (
	"os"
	"os/exec"
)

// Environment represents an environment where the tools are located and can be executed.
type Environment interface {
	// FindExecutable try to find an executable in the environment
	FindExecutable(name string) *Executable

	// RunExecutable run an executable in the environment
	RunExecutable(path string, args []string) error
}

// SystemEnvironment is the operating system environment
var SystemEnvironment Environment = &systemEnvironment{}

type systemEnvironment struct{}

func (s *systemEnvironment) FindExecutable(name string) *Executable {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil
	}

	return &Executable{Path: path, Environment: s}
}

func (s *systemEnvironment) RunExecutable(path string, args []string) error {
	command := exec.Command(path, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}
