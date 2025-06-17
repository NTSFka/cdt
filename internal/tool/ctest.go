package tool

import (
	. "cdt/internal"
)

// CTest is a test runner from CMake
type CTest struct {
	ExecutableTool
}

// NewCTest creates a ctest tool from a custom executable
func NewCTest(executable *Executable) *CTest {
	return &CTest{
		MakeExecutableTool(
			"ctest",
			"CTest",
			"The ctest executable is the CMake test driver program.",
			executable,
		),
	}
}

// DetectCTest create ctest tool can be used in the project
func DetectCTest(environment Environment) *CTest {
	return NewCTest(environment.FindExecutable("ctest"))
}

func (c *CTest) Run(project Project, args []string) error {
	return c.ExecutableTool.Run(project, append(args,
		"--test-dir", project.BuildDirectory(),
	))
}
