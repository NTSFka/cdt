package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCTest_DetectCTest(t *testing.T) {
	ctest := test.RunDetectToolLastFound(t, ctestDetect, []string{"ctest"})

	assert.Equal(t, tool.IdCTest, ctest.Id())
}

func TestCTest_DetectCTest_NotFound(t *testing.T) {
	ctest := test.RunDetectToolNotFound(t, ctestDetect, []string{"ctest"})

	assert.Equal(t, tool.IdCTest, ctest.Id())
}

func TestCTest_DetectCTest_Config(t *testing.T) {
	ctest := test.RunDetectToolConfig(t, ctestDetect, "ctest")

	assert.Equal(t, tool.IdCTest, ctest.Id())
}

func TestCTest_RunForProject_NoOutputDirectory(t *testing.T) {
	test.InvokeRunForProject(
		t,
		ctestTool,
		internal.ProjectInfo{Directory: "project", OutputDirectory: nil},
		[]string{},
		func(exec *test.Executable, executablePath string) {},
		func(err error) {
			require.ErrorIs(t, err, internal.ErrNoOutputDirectory)
		},
	)
}

func TestCTest_RunForProject_Success(t *testing.T) {
	outputDir := "build"

	test.InvokeRunForProjectSuccess(
		t,
		ctestTool,
		internal.ProjectInfo{Directory: "project", OutputDirectory: &outputDir},
		[]string{},
		[]string{"--test-dir", outputDir},
	)
}

func TestCTest_RunForProject_Failed(t *testing.T) {
	outputDir := "build"

	test.InvokeRunForProjectFailed(
		t,
		ctestTool,
		internal.ProjectInfo{Directory: "project", OutputDirectory: &outputDir},
		[]string{},
		[]string{"--test-dir", outputDir},
	)
}

func ctestDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectCTest(ctx, options)
}

func ctestTool(executable func() (*internal.Executable, error)) internal.Tool {
	return tool.NewCTest(executable)
}
