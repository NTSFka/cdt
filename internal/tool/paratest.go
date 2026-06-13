package tool

import (
	"context"
	"path/filepath"

	"cdt/internal"
)

const IdParaTest = "paratest"

type ParaTest struct {
	// Paratest is just PHPUnit with parallel testing support.
	PHPUnit
}

// DetectParaTest create a tool for paratest.
func DetectParaTest(
	ctx context.Context,
	options internal.DetectOptions,
) *ParaTest {
	if path, ok := options.ToolsPaths[IdParaTest]; ok {
		return NewParaTest(func() (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, path)
		})
	}

	return NewParaTest(internal.DetectExecutableChain(
		[]string{filepath.Join(options.ProjectDirectory, "vendor/bin/paratest"), "paratest"},
		func(name string) (*internal.Executable, error) {
			return options.Environment.FindExecutable(ctx, name)
		},
	))
}

// NewParaTest creates a paratest tool from a custom executable.
func NewParaTest(detect internal.ExecutableToolDetectFunc) *ParaTest {
	return &ParaTest{
		PHPUnit{
			ExecutableTool: internal.MakeExecutableTool(
				IdParaTest,
				"ParaTest",
				"The objective of ParaTest is to support parallel testing in PHPUnit.",
				internal.Tags{internal.ToolTagPhp, internal.ToolTagTest},
				detect,
			),
		},
	}
}
