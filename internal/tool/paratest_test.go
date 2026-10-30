package tool_test

import (
	"context"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestParaTest_DetectParaTest_Composer(t *testing.T) {
	paraTest := test.RunDetectToolLastFound(
		t,
		paraTestDetect,
		[]string{filepath.Join("vendor", "bin", "paratest")},
	)

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_DetectParaTest_System(t *testing.T) {
	paraTest := test.RunDetectToolLastFound(
		t,
		paraTestDetect,
		[]string{filepath.Join("vendor", "bin", "paratest"), "paratest"},
	)

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_DetectParaTest_NotFound(t *testing.T) {
	paraTest := test.RunDetectToolNotFound(
		t,
		paraTestDetect,
		[]string{filepath.Join("vendor", "bin", "paratest"), "paratest"},
	)

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_DetectParaTest_Config(t *testing.T) {
	paraTest := test.RunDetectToolConfig(t, paraTestDetect, "paratest")

	assert.Equal(t, tool.IdParaTest, paraTest.Id())
}

func TestParaTest_ParaTest_RunTests(t *testing.T) {
	test.InvokeRunTestsSuccess(
		t,
		paraTestTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{},
	)
}

func TestParaTest_ParaTest_RunTests_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeRunTestsSuccess(
		t,
		paraTestTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{pattern},
	)
}

func TestParaTest_ParaTest_CollectCoverage(t *testing.T) {
	test.InvokeCollectCoverageSuccess(
		t,
		paraTestCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{"--coverage-text"},
	)
}

func TestParaTest_ParaTest_CollectCoverage_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeCollectCoverageSuccess(
		t,
		paraTestCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{"--coverage-text", pattern},
	)
}

func paraTestDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectParaTest(ctx, options)
}

func paraTestTester(executable func() (*internal.Executable, error)) internal.ProjectTester {
	return tool.NewParaTest(executable)
}

func paraTestCoverageCollector(
	executable func() (*internal.Executable, error),
) internal.ProjectCoverageCollector {
	return tool.NewParaTest(executable)
}
