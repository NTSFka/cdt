package tool_test

import (
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPyTest_DetectPyTest(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(env.NewExecutable("/bin/pytest"), nil)

	pyTest := tool.DetectPyTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, pyTest)
	assert.Equal(t, "pytest", pyTest.Id())
	assert.True(t, pyTest.IsAvailable())

	if executable := pyTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pytest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPyTest_DetectPyTest_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest").
		Return(nil, nil)

	pyTest := tool.DetectPyTest(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, pyTest)
	assert.Equal(t, "pytest", pyTest.Id())
	assert.False(t, pyTest.IsAvailable())
	assert.Nil(t, pyTest.Executable())

	env.AssertExpectations(t)
}

func TestPyTest_DetectPyTest_Config(t *testing.T) {
	env := test.NewEnvironment(t)

	env.OnFindExecutable("pytest2").
		Return(env.NewExecutable("/bin/pytest"), nil)

	pyTest := tool.DetectPyTest(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"pytest": "pytest2"},
	})
	assert.NotNil(t, pyTest)
	assert.Equal(t, "pytest", pyTest.Id())
	assert.True(t, pyTest.IsAvailable())

	if executable := pyTest.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/pytest", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestPyTest_PyTest_RunTests(t *testing.T) {
	exec := test.NewExecutable(t)

	pyTest := tool.NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{}).
		Return(nil)

	err := pyTest.RunTests(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPyTest_PyTest_RunTests_Pattern(t *testing.T) {
	exec := test.NewExecutable(t)

	pyTest := tool.NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("test", []string{"tests/*"}).
		Return(nil)

	err := pyTest.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: info, Pattern: internal.StrPtr("tests/*")},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestPyTest_RunTests_Pattern_FormatUnsupported(t *testing.T) {
	exec := test.NewExecutable(t)

	pyTest := tool.NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}

	data := []internal.TestsReportFormat{
		"test",
		internal.TestsReportFormatRawEvents,
		internal.TestsReportFormatJson,
		internal.TestsReportFormatCtrf,
		internal.TestsReportFormatTeamCity,
	}

	for _, format := range data {
		t.Run(string(format), func(t *testing.T) {
			err := pyTest.RunTests(
				t.Context(),
				internal.ProjectTesterOptions{
					ProjectInfo: info,
					Output:      internal.OutputOptions[internal.TestsReportFormat]{Format: format},
					Pattern:     internal.StrPtr("test1"),
				},
			)
			require.EqualError(t, err, "unsupported report format: "+string(format))

			exec.AssertExpectations(t)
		})
	}
}

func TestPyTest_RunTests_Pattern_FormatJUnit(t *testing.T) {
	exec := test.NewExecutable(t)

	pyTest := tool.NewPyTest(exec.LazyExecutable("test"))

	info := internal.ProjectInfo{Directory: "."}
	reportFilename := "pytest-junit.xml"

	exec.OnRun("test", []string{"test1", "--junitxml=" + reportFilename}).
		Return(nil)

	err := pyTest.RunTests(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatJUnit,
				Filename: &reportFilename,
			},
			Pattern: internal.StrPtr("test1"),
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}
