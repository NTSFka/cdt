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

func TestPHPUnit_DetectPHPUnit_Composer(t *testing.T) {
	phpUnit := test.RunDetectToolLastFound(
		t,
		phpUnitDetect,
		[]string{filepath.Join("vendor", "bin", "phpunit")},
	)

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_DetectPHPUnit_System(t *testing.T) {
	phpUnit := test.RunDetectToolLastFound(
		t,
		phpUnitDetect,
		[]string{filepath.Join("vendor", "bin", "phpunit"), "phpunit"},
	)

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_DetectPHPUnit_NotFound(t *testing.T) {
	phpUnit := test.RunDetectToolNotFound(
		t,
		phpUnitDetect,
		[]string{filepath.Join("vendor", "bin", "phpunit"), "phpunit"},
	)

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_DetectPHPUnit_Config(t *testing.T) {
	phpUnit := test.RunDetectToolConfig(t, phpUnitDetect, "phpunit")

	assert.Equal(t, tool.IdPHPUnit, phpUnit.Id())
}

func TestPHPUnit_RunTests(t *testing.T) {
	test.InvokeRunTestsSuccess(
		t,
		phpUnitTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{},
	)
}

func TestPHPUnit_RunTests_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeRunTestsSuccess(
		t,
		phpUnitTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{pattern},
	)
}

func TestPHPUnit_CollectCoverage(t *testing.T) {
	test.InvokeCollectCoverageSuccess(
		t,
		phpUnitCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
		},
		[]string{"--coverage-text"},
	)
}

func TestPHPUnit_CollectCoverage_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeCollectCoverageSuccess(
		t,
		phpUnitCoverageCollector,
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{"--coverage-text", pattern},
	)
}

func phpUnitDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPHPUnit(ctx, options)
}

func phpUnitTester(executable func() (*internal.Executable, error)) internal.ProjectTester {
	return tool.NewPHPUnit(executable)
}

func phpUnitCoverageCollector(
	executable func() (*internal.Executable, error),
) internal.ProjectCoverageCollector {
	return tool.NewPHPUnit(executable)
}
