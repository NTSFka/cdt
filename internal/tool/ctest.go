package tool

import (
	. "cdt/internal"
	"errors"
)

// CTest is a test runner from CMake
type CTest struct {
	executable *Executable
}

// NewCTest creates a ctest tool from a custom executable
func NewCTest(executable *Executable) *CTest {
	return &CTest{
		executable: executable,
	}
}

// DetectCTest create ctest tool can be used in the project
func DetectCTest() *CTest {
	return NewCTest(FindExecutable("ctest"))
}

func (c *CTest) Id() string {
	return "ctest"
}

func (c *CTest) Name() string {
	return "CTest"
}

func (c *CTest) Info() string {
	return "The ctest executable is the CMake test driver program."
}

func (c *CTest) ExecutablePath() *string {
	if c.executable != nil {
		return &c.executable.Path
	}

	return nil
}

func (c *CTest) IsAvailable() bool {
	return c.executable != nil
}

func (c *CTest) Run(project Project, args []string) error {
	if c.executable == nil {
		return errors.New("CTest is not installed on the system")
	}

	callArgs := args
	callArgs = append(callArgs, "--test-dir", project.BuildDirectory())

	return c.executable.Run(callArgs)
}
