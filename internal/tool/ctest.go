package tool

import (
	"context"

	"cdt/internal"
)

const IdCTest = "ctest"

// CTest is a test runner from CMake.
type CTest struct {
	internal.ExecutableTool
}

// NewCTest creates a ctest tool from a custom executable.
func NewCTest(detect internal.ExecutableToolDetectFunc) *CTest {
	return &CTest{
		internal.MakeExecutableTool(
			IdCTest,
			"CTest",
			"The ctest executable is the CMake test driver program.",
			internal.Tags{internal.ToolTagC, internal.ToolTagCpp, internal.ToolTagTest},
			detect,
		),
	}
}

// DetectCTest create ctest tool can be used in the project.
func DetectCTest(
	ctx context.Context,
	config internal.ConfigTools,
	environment internal.Environment,
) *CTest {
	return NewCTest(func() (*internal.Executable, error) {
		return environment.FindExecutable(ctx, config.Get(IdCTest, "ctest"))
	})
}

func (c *CTest) RunForProject(ctx context.Context, info internal.ProjectInfo, args []string) error {
	if info.OutputDirectory == nil {
		return internal.ErrNoOutputDirectory
	}

	return c.ExecutableTool.RunForProject(ctx, info, append(args,
		"--test-dir", *info.OutputDirectory,
	))
}
