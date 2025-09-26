package workflow_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/workflow"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createConfiguratorTool(id string, executable *internal.Executable) *struct {
	internal.ExecutableTool
	test.ProjectConfigurator
} {
	return &struct {
		internal.ExecutableTool
		test.ProjectConfigurator
	}{
		internal.MakeExecutableTool(id, "TestPattern", "", internal.Tags{}, func() *internal.Executable {
			return executable
		}),
		test.ProjectConfigurator{},
	}
}

func TestConfiguratorFallback_Details_Empty(t *testing.T) {
	fallback := &workflow.ConfiguratorFallback{}

	assert.Equal(t, "fallback ()", fallback.Details())
}

func TestConfiguratorFallback_Details(t *testing.T) {
	tool1 := createConfiguratorTool("test1", nil)
	tool2 := createConfiguratorTool("test2", nil)

	fallback := &workflow.ConfiguratorFallback{tool1, tool2}

	assert.Equal(t, "fallback (test1, test2)", fallback.Details())
}

func TestConfiguratorFallback_Configure_Empty(t *testing.T) {
	fallback := &workflow.ConfiguratorFallback{}

	err := fallback.Configure(t.Context(), internal.ProjectConfiguratorOptions{})

	require.EqualError(t, err, "no configurator tool available: none")
}

func TestConfiguratorFallback_Configure_NoAvailable(t *testing.T) {
	tool1 := createConfiguratorTool("test1", nil)
	tool2 := createConfiguratorTool("test2", nil)

	fallback := &workflow.ConfiguratorFallback{tool1, tool2}

	err := fallback.Configure(t.Context(), internal.ProjectConfiguratorOptions{})

	require.EqualError(t, err, "no configurator tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestConfiguratorFallback_Configure_Available1(t *testing.T) {
	tool1 := createConfiguratorTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createConfiguratorTool("test2", nil)

	fallback := &workflow.ConfiguratorFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("Configure", mock.Anything, mock.Anything).Return(nil)

	err := fallback.Configure(t.Context(), internal.ProjectConfiguratorOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestConfiguratorFallback_Configure_Available2(t *testing.T) {
	tool1 := createConfiguratorTool("test1", nil)
	tool2 := createConfiguratorTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.ConfiguratorFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("Configure", mock.Anything, mock.Anything).Return(nil)

	err := fallback.Configure(t.Context(), internal.ProjectConfiguratorOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func createBuilderTool(id string, executable *internal.Executable) *struct {
	internal.ExecutableTool
	test.ProjectBuilder
} {
	return &struct {
		internal.ExecutableTool
		test.ProjectBuilder
	}{
		internal.MakeExecutableTool(id, "Test", "", internal.Tags{}, func() *internal.Executable {
			return executable
		}),
		test.ProjectBuilder{},
	}
}

func TestBuilderFallback_Details_Empty(t *testing.T) {
	fallback := &workflow.BuilderFallback{}

	assert.Equal(t, "fallback ()", fallback.Details())
}

func TestBuilderFallback_Details(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", nil)

	fallback := &workflow.BuilderFallback{tool1, tool2}

	assert.Equal(t, "fallback (test1, test2)", fallback.Details())
}

func TestBuilderFallback_BuildAll_Empty(t *testing.T) {
	fallback := &workflow.BuilderFallback{}

	err := fallback.BuildAll(t.Context(), internal.ProjectBuilderOptions{})

	require.EqualError(t, err, "no builder tool available: none")
}

func TestBuilderFallback_BuildAll_NoAvailable(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", nil)

	fallback := &workflow.BuilderFallback{tool1, tool2}

	err := fallback.BuildAll(t.Context(), internal.ProjectBuilderOptions{})

	require.EqualError(t, err, "no builder tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildAll_Available1(t *testing.T) {
	tool1 := createBuilderTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createBuilderTool("test2", nil)

	fallback := &workflow.BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("BuildAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.BuildAll(t.Context(), internal.ProjectBuilderOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildAll_Available2(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("BuildAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.BuildAll(t.Context(), internal.ProjectBuilderOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildTargets_Empty(t *testing.T) {
	fallback := &workflow.BuilderFallback{}

	err := fallback.BuildTargets(t.Context(), internal.ProjectBuilderOptions{}, []string{"target1"})

	require.EqualError(t, err, "no builder tool available: none")
}

func TestBuilderFallback_BuildTargets_NoAvailable(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", nil)

	fallback := &workflow.BuilderFallback{tool1, tool2}

	err := fallback.BuildTargets(t.Context(), internal.ProjectBuilderOptions{}, []string{"target1"})

	require.EqualError(t, err, "no builder tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildTargets_Available1(t *testing.T) {
	tool1 := createBuilderTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createBuilderTool("test2", nil)

	fallback := &workflow.BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("BuildTargets", mock.Anything, mock.Anything, []string{"target1"}).Return(nil)

	err := fallback.BuildTargets(t.Context(), internal.ProjectBuilderOptions{}, []string{"target1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildTargets_Available2(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("BuildTargets", mock.Anything, mock.Anything, []string{"target1"}).Return(nil)

	err := fallback.BuildTargets(t.Context(), internal.ProjectBuilderOptions{}, []string{"target1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func createTesterTool(id string, executable *internal.Executable) *struct {
	internal.ExecutableTool
	test.ProjectTester
} {
	return &struct {
		internal.ExecutableTool
		test.ProjectTester
	}{
		internal.MakeExecutableTool(id, "Test", "", internal.Tags{}, func() *internal.Executable {
			return executable
		}),
		test.ProjectTester{},
	}
}

func TestTesterFallback_Details_Empty(t *testing.T) {
	fallback := &workflow.TesterFallback{}

	assert.Equal(t, "fallback ()", fallback.Details())
}

func TestTesterFallback_Details(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", nil)

	fallback := &workflow.TesterFallback{tool1, tool2}

	assert.Equal(t, "fallback (test1, test2)", fallback.Details())
}

func TestTesterFallback_TestAll_Empty(t *testing.T) {
	fallback := &workflow.TesterFallback{}

	err := fallback.TestAll(t.Context(), internal.ProjectTesterOptions{})

	require.EqualError(t, err, "no tester tool available: none")
}

func TestTesterFallback_TestAll_NoAvailable(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", nil)

	fallback := &workflow.TesterFallback{tool1, tool2}

	err := fallback.TestAll(t.Context(), internal.ProjectTesterOptions{})

	require.EqualError(t, err, "no tester tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_TestAll_Available1(t *testing.T) {
	tool1 := createTesterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createTesterTool("test2", nil)

	fallback := &workflow.TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("TestAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.TestAll(t.Context(), internal.ProjectTesterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_TestAll_Available2(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("TestAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.TestAll(t.Context(), internal.ProjectTesterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_Test_Empty(t *testing.T) {
	fallback := &workflow.TesterFallback{}

	err := fallback.TestPattern(t.Context(), internal.ProjectTesterOptions{}, "target1")

	require.EqualError(t, err, "no tester tool available: none")
}

func TestTesterFallback_Test_NoAvailable(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", nil)

	fallback := &workflow.TesterFallback{tool1, tool2}

	err := fallback.TestPattern(t.Context(), internal.ProjectTesterOptions{}, "target1")

	require.EqualError(t, err, "no tester tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_Test_Available1(t *testing.T) {
	tool1 := createTesterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createTesterTool("test2", nil)

	fallback := &workflow.TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("TestPattern", mock.Anything, mock.Anything, "target1").Return(nil)

	err := fallback.TestPattern(t.Context(), internal.ProjectTesterOptions{}, "target1")

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_Test_Available2(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("TestPattern", mock.Anything, mock.Anything, "target1").Return(nil)

	err := fallback.TestPattern(t.Context(), internal.ProjectTesterOptions{}, "target1")

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func createFormatterTool(id string, executable *internal.Executable) *struct {
	internal.ExecutableTool
	test.ProjectFormatter
} {
	return &struct {
		internal.ExecutableTool
		test.ProjectFormatter
	}{
		internal.MakeExecutableTool(id, "Test", "", internal.Tags{}, func() *internal.Executable {
			return executable
		}),
		test.ProjectFormatter{},
	}
}

func TestFormatterFallback_Details_Empty(t *testing.T) {
	fallback := &workflow.FormatterFallback{}

	assert.Equal(t, "fallback ()", fallback.Details())
}

func TestFormatterFallback_Details(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	assert.Equal(t, "fallback (test1, test2)", fallback.Details())
}

func TestFormatterFallback_FormatAll_Empty(t *testing.T) {
	fallback := &workflow.FormatterFallback{}

	err := fallback.FormatAll(t.Context(), internal.ProjectFormatterOptions{})

	require.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatAll_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	err := fallback.FormatAll(t.Context(), internal.ProjectFormatterOptions{})

	require.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatAll_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatAll(t.Context(), internal.ProjectFormatterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatAll_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatAll(t.Context(), internal.ProjectFormatterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatFiles_Empty(t *testing.T) {
	fallback := &workflow.FormatterFallback{}

	err := fallback.FormatFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatFiles_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	err := fallback.FormatFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatFiles_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatFiles_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckAll_Empty(t *testing.T) {
	fallback := &workflow.FormatterFallback{}

	err := fallback.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{})

	require.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatCheckAll_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	err := fallback.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{})

	require.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckAll_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatCheckAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckAll_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatCheckAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckFiles_Empty(t *testing.T) {
	fallback := &workflow.FormatterFallback{}

	err := fallback.FormatCheckFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatCheckFiles_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	err := fallback.FormatCheckFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckFiles_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatCheckFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatCheckFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckFiles_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatCheckFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatCheckFiles(t.Context(), internal.ProjectFormatterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func createLinterTool(id string, executable *internal.Executable) *struct {
	internal.ExecutableTool
	test.ProjectLinter
} {
	return &struct {
		internal.ExecutableTool
		test.ProjectLinter
	}{
		internal.MakeExecutableTool(id, "Test", "", internal.Tags{}, func() *internal.Executable {
			return executable
		}),
		test.ProjectLinter{},
	}
}

func TestLinterFallback_Details_Empty(t *testing.T) {
	fallback := &workflow.LinterFallback{}

	assert.Equal(t, "fallback ()", fallback.Details())
}

func TestLinterFallback_Details(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	fallback := &workflow.LinterFallback{tool1, tool2}

	assert.Equal(t, "fallback (test1, test2)", fallback.Details())
}

func TestLinterFallback_LintAll_Empty(t *testing.T) {
	fallback := &workflow.LinterFallback{}

	err := fallback.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.EqualError(t, err, "no linter tool available: none")
}

func TestLinterFallback_LintAll_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	fallback := &workflow.LinterFallback{tool1, tool2}

	err := fallback.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintAll_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	fallback := &workflow.LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintAll_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintFiles_Empty(t *testing.T) {
	fallback := &workflow.LinterFallback{}

	err := fallback.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no linter tool available: none")
}

func TestLinterFallback_LintFiles_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	fallback := &workflow.LinterFallback{tool1, tool2}

	err := fallback.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintFiles_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	fallback := &workflow.LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintFiles_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_Details_Empty(t *testing.T) {
	list := &workflow.LinterList{}

	assert.Equal(t, "list ()", list.Details())
}

func TestLinterList_Details(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	list := &workflow.LinterList{tool1, tool2}

	assert.Equal(t, "list (test1, test2)", list.Details())
}

func TestLinterList_LintAll_Empty(t *testing.T) {
	list := &workflow.LinterList{}

	err := list.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.EqualError(t, err, "no linter tool available: none")
}

func TestLinterList_LintAll_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	list := &workflow.LinterList{tool1, tool2}

	err := list.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintAll_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	list := &workflow.LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := list.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintAll_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	list := &workflow.LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintAll", mock.Anything, mock.Anything).Return(nil)
	tool2.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := list.LintAll(t.Context(), internal.ProjectLinterOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintFiles_Empty(t *testing.T) {
	list := &workflow.LinterList{}

	err := list.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no linter tool available: none")
}

func TestLinterList_LintFiles_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	list := &workflow.LinterList{tool1, tool2}

	err := list.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintFiles_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	list := &workflow.LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := list.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintFiles_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	list := &workflow.LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)
	tool2.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := list.LintFiles(t.Context(), internal.ProjectLinterOptions{}, []string{"file1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func createRunnerTool(id string, executable *internal.Executable) *struct {
	internal.ExecutableTool
	test.ProjectRunner
} {
	return &struct {
		internal.ExecutableTool
		test.ProjectRunner
	}{
		internal.MakeExecutableTool(id, "Test", "", internal.Tags{}, func() *internal.Executable {
			return executable
		}),
		test.ProjectRunner{},
	}
}

func TestRunnerFallback_Details_Empty(t *testing.T) {
	fallback := &workflow.RunnerFallback{}

	assert.Equal(t, "fallback ()", fallback.Details())
}

func TestRunnerFallback_Details(t *testing.T) {
	tool1 := createRunnerTool("test1", nil)
	tool2 := createRunnerTool("test2", nil)

	fallback := &workflow.RunnerFallback{tool1, tool2}

	assert.Equal(t, "fallback (test1, test2)", fallback.Details())
}

func TestRunnerFallback_RunTarget_Empty(t *testing.T) {
	fallback := &workflow.RunnerFallback{}

	err := fallback.RunTarget(t.Context(), internal.ProjectRunnerOptions{}, "target1")

	require.EqualError(t, err, "no runner tool available: none")
}

func TestRunnerFallback_RunTarget_NoAvailable(t *testing.T) {
	tool1 := createRunnerTool("test1", nil)
	tool2 := createRunnerTool("test2", nil)

	fallback := &workflow.RunnerFallback{tool1, tool2}

	err := fallback.RunTarget(t.Context(), internal.ProjectRunnerOptions{}, "target1")

	require.EqualError(t, err, "no runner tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestRunnerFallback_RunTarget_Available1(t *testing.T) {
	tool1 := createRunnerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createRunnerTool("test2", nil)

	fallback := &workflow.RunnerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("RunTarget", mock.Anything, mock.Anything, "target1").Return(nil)

	err := fallback.RunTarget(t.Context(), internal.ProjectRunnerOptions{}, "target1")

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestRunnerFallback_RunTarget_Available2(t *testing.T) {
	tool1 := createRunnerTool("test1", nil)
	tool2 := createRunnerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.RunnerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("RunTarget", mock.Anything, mock.Anything, "target1").Return(nil)

	err := fallback.RunTarget(t.Context(), internal.ProjectRunnerOptions{}, "target1")

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func createDependencyManagerTool(id string, executable *internal.Executable) *struct {
	internal.ExecutableTool
	test.DependencyManager
} {
	return &struct {
		internal.ExecutableTool
		test.DependencyManager
	}{
		internal.MakeExecutableTool(id, "Test", "", internal.Tags{}, func() *internal.Executable {
			return executable
		}),
		test.DependencyManager{},
	}
}

func TestDependencyManagerFallback_Details_Empty(t *testing.T) {
	fallback := &workflow.DependencyManagerFallback{}

	assert.Equal(t, "fallback ()", fallback.Details())
}

func TestDependencyManagerFallback_Details(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	assert.Equal(t, "fallback (test1, test2)", fallback.Details())
}

func TestDependencyManagerFallback_AddDependencies_Empty(t *testing.T) {
	fallback := &workflow.DependencyManagerFallback{}

	err := fallback.AddDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_AddDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	err := fallback.AddDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AddDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("AddDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.AddDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AddDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("AddDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.AddDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_RemoveDependencies_Empty(t *testing.T) {
	fallback := &workflow.DependencyManagerFallback{}

	err := fallback.RemoveDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_RemoveDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	err := fallback.RemoveDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_RemoveDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("RemoveDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.RemoveDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_RemoveDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("RemoveDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.RemoveDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"}, false)

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_UpdateDependencies_Empty(t *testing.T) {
	fallback := &workflow.DependencyManagerFallback{}

	err := fallback.UpdateDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"})

	require.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_UpdateDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	err := fallback.UpdateDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"})

	require.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_UpdateDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("UpdateDependencies", mock.Anything, mock.Anything, []string{"dep1"}).Return(nil)

	err := fallback.UpdateDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_UpdateDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("UpdateDependencies", mock.Anything, mock.Anything, []string{"dep1"}).Return(nil)

	err := fallback.UpdateDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, []string{"dep1"})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_FetchDependencies_Empty(t *testing.T) {
	fallback := &workflow.DependencyManagerFallback{}

	err := fallback.FetchDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, false)

	require.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_FetchDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	err := fallback.FetchDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, false)

	require.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_FetchDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FetchDependencies", mock.Anything, mock.Anything, false).Return(nil)

	err := fallback.FetchDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, false)

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_FetchDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FetchDependencies", mock.Anything, mock.Anything, false).Return(nil)

	err := fallback.FetchDependencies(t.Context(), internal.ProjectDependencyManagerOptions{}, false)

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_ListDependencies_Empty(t *testing.T) {
	fallback := &workflow.DependencyManagerFallback{}

	err := fallback.ListDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_ListDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	err := fallback.ListDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_ListDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("ListDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.ListDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_ListDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("ListDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.ListDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AuditDependencies_Empty(t *testing.T) {
	fallback := &workflow.DependencyManagerFallback{}

	err := fallback.AuditDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_AuditDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	err := fallback.AuditDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AuditDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("AuditDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.AuditDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AuditDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &workflow.DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("AuditDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.AuditDependencies(t.Context(), internal.ProjectDependencyManagerOptions{})

	require.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}
