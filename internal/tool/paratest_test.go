package tool_test

import (
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParaTest_DetectParaTest_Composer(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "paratest")).
		Return(env.NewExecutable("/bin/paratest"), nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.True(t, paraTest.IsAvailable())

	if executable := paraTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/paratest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_System(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable(filepath.Join("vendor", "bin", "paratest")).
		Return(nil, nil)

	// System installation
	env.OnFindExecutable("paratest").
		Return(env.NewExecutable("/bin/paratest"), nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.True(t, paraTest.IsAvailable())

	if executable := paraTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/paratest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable(filepath.Join("vendor", "bin", "paratest")).
		Return(nil, nil)
	env.OnFindExecutable("paratest").
		Return(nil, nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.False(t, paraTest.IsAvailable())
	assert.Nil(t, paraTest.Executable())

	env.AssertExpectations(t)
}

func TestParaTest_DetectParaTest_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	// Composer installation
	env.OnFindExecutable("paratest-7").
		Return(env.NewExecutable("/bin/paratest"), nil)

	paraTest := tool.DetectParaTest(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"paratest": "paratest-7"},
	})
	assert.NotNil(t, paraTest)
	assert.Equal(t, "paratest", paraTest.Id())
	assert.True(t, paraTest.IsAvailable())

	if executable := paraTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/paratest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestParaTest_ParaTest_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := paraTest.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_ParaTest_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := paraTest.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: info},
		"tests/*",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_TestPattern_FormatUnsupported(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	data := []internal.TestsReportFormat{
		"test",
		internal.TestsReportFormatJson,
		internal.TestsReportFormatCtrf,
	}

	for _, format := range data {
		t.Run(string(format), func(t *testing.T) {
			err := paraTest.TestPattern(
				t.Context(),
				internal.ProjectTesterOptions{
					ProjectInfo: info,
					Output:      internal.OutputOptions[internal.TestsReportFormat]{Format: format},
				},
				"test1",
			)
			require.EqualError(t, err, "unsupported report format: "+string(format))

			exec.AssertExpectations(t)
		})
	}
}

func TestParaTest_TestPattern_FormatRawEvents(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	reportFilename := "paratest-raw-events.txt"

	exec.OnRun("test", []string{"test1", "--log-events-text", reportFilename}).
		Return(nil)

	err := paraTest.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatRawEvents,
				Filename: &reportFilename,
			},
		},
		"test1",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_TestPattern_FormatJUnit(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	reportFilename := "paratest-junit.xml"

	exec.OnRun("test", []string{"test1", "--log-junit", reportFilename}).
		Return(nil)

	err := paraTest.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatJUnit,
				Filename: &reportFilename,
			},
		},
		"test1",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_TestPattern_FormatTeamCity(t *testing.T) {
	exec := test.NewExecutable(t)

	paraTest := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	reportFilename := "paratest-teamcity.xml"

	exec.OnRun("test", []string{"test1", "--log-teamcity", reportFilename}).
		Return(nil)

	err := paraTest.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatTeamCity,
				Filename: &reportFilename,
			},
		},
		"test1",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_ParaTest_CollectCoverage(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"--coverage-text"}).
		Return(nil)

	err := phpUnit.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestParaTest_ParaTest_CollectCoverage_Pattern(t *testing.T) {
	exec := test.NewExecutable(t)

	phpUnit := tool.NewParaTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"--coverage-text", "tests/*"}).
		Return(nil)

	err := phpUnit.CollectCoverage(
		t.Context(),
		internal.ProjectCoverageCollectorOptions{
			ProjectInfo: info,
			Pattern:     internal.StrPtr("tests/*"),
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
