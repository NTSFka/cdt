package tool

import (
	"cdt/internal"
	"context"
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
			internal.Tags{internal.ToolTagC, internal.ToolTagCpp, internal.ToolTagTest},
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

func (c *CTest) RunForProject(ctx context.Context, info internal.ProjectInfo, args []string) error {
	if info.IntermediateDirectory == nil {
		return internal.ErrNoIntermediateDirectory
	}

	return c.ExecutableTool.RunForProject(ctx, info, append(args,
		"--test-dir", *info.IntermediateDirectory,
	))
}
