package tool_test

import (
	"context"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
)

func TestPyTest_DetectPyTest(t *testing.T) {
	pyTest := test.RunDetectToolLastFound(t, pyTestDetect, []string{"pytest"})

	assert.Equal(t, tool.IdPyTest, pyTest.Id())
}

func TestPyTest_DetectPyTest_NotFound(t *testing.T) {
	pyTest := test.RunDetectToolNotFound(t, pyTestDetect, []string{"pytest"})

	assert.Equal(t, tool.IdPyTest, pyTest.Id())
}

func TestPyTest_DetectPyTest_Config(t *testing.T) {
	pyTest := test.RunDetectToolConfig(t, pyTestDetect, "pytest")

	assert.Equal(t, tool.IdPyTest, pyTest.Id())
}

func TestPyTest_PyTest_RunTests(t *testing.T) {
	test.InvokeRunTestsSuccess(
		t,
		pyTestTester,
		internal.ProjectTesterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{},
	)
}

func TestPyTest_PyTest_RunTests_Pattern(t *testing.T) {
	pattern := "tests/*"

	test.InvokeRunTestsSuccess(
		t,
		pyTestTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     &pattern,
		},
		[]string{pattern},
	)
}

func TestPyTest_RunTests_Pattern_Format_Unsupported(t *testing.T) {
	dataSet := []internal.TestsReportFormat{
		"test",
		internal.TestsReportFormatRawEvents,
		internal.TestsReportFormatJson,
		internal.TestsReportFormatCtrf,
		internal.TestsReportFormatTeamCity,
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeRunTestsOutputFormatUnsupported(
				t,
				pyTestTester,
				format,
				internal.StrPtr("test1"),
			)
		})
	}
}

func TestPyTest_RunTests_Pattern_Format_Supported(t *testing.T) {
	dataSet := []struct {
		Format   internal.TestsReportFormat
		Filename *string
		Args     []string
	}{
		{
			Format:   internal.TestsReportFormatJUnit,
			Filename: internal.StrPtr("pytest-junit.xml"),
			Args:     []string{"--junitxml=pytest-junit.xml"},
		},
	}

	for _, values := range dataSet {
		t.Run(string(values.Format), func(t *testing.T) {
			test.InvokeRunTestsSuccess(
				t,
				pyTestTester,
				internal.ProjectTesterOptions{
					ProjectInfo: internal.ProjectInfo{Directory: "."},
					Output: internal.OutputOptions[internal.TestsReportFormat]{
						Format:   values.Format,
						Filename: values.Filename,
					},
				},
				values.Args,
			)
		})
	}
}

func pyTestDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectPyTest(ctx, options)
}

func pyTestTester(executable func() (*internal.Executable, error)) internal.ProjectTester {
	return tool.NewPyTest(executable)
}
