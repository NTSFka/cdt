package tool

import (
	"cdt/internal"
)

// CTest is a test runner from CMake
type CTest struct {
	internal.ExecutableTool
}

// NewCTest creates a ctest tool from a custom executable
func NewCTest(detect func() *internal.Executable) *CTest {
	return &CTest{
		internal.MakeExecutableTool(
			"ctest",
			"CTest",
			"The ctest executable is the CMake test driver program.",
			detect,
		),
	}
}

// DetectCTest create ctest tool can be used in the project
func DetectCTest(environment internal.Environment) *CTest {
	return NewCTest(func() *internal.Executable {
		return environment.FindExecutable("ctest")
	})
}

func (c *CTest) Run(project internal.Project, args []string) error {
	return c.ExecutableTool.Run(project, append(args,
		"--test-dir", project.BuildDirectory(),
	))
}
