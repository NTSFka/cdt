package project

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestConfiguratorFallback_Configure_Empty(t *testing.T) {
	fallback := &ConfiguratorFallback{}

	err := fallback.Configure(context.Background(), internal.ProjectConfiguratorOptions{})

	assert.EqualError(t, err, "no configurator tool available: none")
}

func TestConfiguratorFallback_Configure_NoAvailable(t *testing.T) {
	tool1 := createConfiguratorTool("test1", nil)
	tool2 := createConfiguratorTool("test2", nil)

	fallback := &ConfiguratorFallback{tool1, tool2}

	err := fallback.Configure(context.Background(), internal.ProjectConfiguratorOptions{})

	assert.EqualError(t, err, "no configurator tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestConfiguratorFallback_Configure_Available1(t *testing.T) {
	tool1 := createConfiguratorTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createConfiguratorTool("test2", nil)

	fallback := &ConfiguratorFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("Configure", mock.Anything, mock.Anything).Return(nil)

	err := fallback.Configure(context.Background(), internal.ProjectConfiguratorOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestConfiguratorFallback_Configure_Available2(t *testing.T) {
	tool1 := createConfiguratorTool("test1", nil)
	tool2 := createConfiguratorTool("test2", &internal.Executable{Path: "test1"})

	fallback := &ConfiguratorFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("Configure", mock.Anything, mock.Anything).Return(nil)

	err := fallback.Configure(context.Background(), internal.ProjectConfiguratorOptions{})

	assert.NoError(t, err)

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

func TestBuilderFallback_BuildAll_Empty(t *testing.T) {
	fallback := &BuilderFallback{}

	err := fallback.BuildAll(context.Background(), internal.ProjectBuilderOptions{})

	assert.EqualError(t, err, "no builder tool available: none")
}

func TestBuilderFallback_BuildAll_NoAvailable(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", nil)

	fallback := &BuilderFallback{tool1, tool2}

	err := fallback.BuildAll(context.Background(), internal.ProjectBuilderOptions{})

	assert.EqualError(t, err, "no builder tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildAll_Available1(t *testing.T) {
	tool1 := createBuilderTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createBuilderTool("test2", nil)

	fallback := &BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("BuildAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.BuildAll(context.Background(), internal.ProjectBuilderOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildAll_Available2(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", &internal.Executable{Path: "test1"})

	fallback := &BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("BuildAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.BuildAll(context.Background(), internal.ProjectBuilderOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildTargets_Empty(t *testing.T) {
	fallback := &BuilderFallback{}

	err := fallback.BuildTargets(context.Background(), internal.ProjectBuilderOptions{}, []string{"target1"})

	assert.EqualError(t, err, "no builder tool available: none")
}

func TestBuilderFallback_BuildTargets_NoAvailable(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", nil)

	fallback := &BuilderFallback{tool1, tool2}

	err := fallback.BuildTargets(context.Background(), internal.ProjectBuilderOptions{}, []string{"target1"})

	assert.EqualError(t, err, "no builder tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildTargets_Available1(t *testing.T) {
	tool1 := createBuilderTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createBuilderTool("test2", nil)

	fallback := &BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("BuildTargets", mock.Anything, mock.Anything, []string{"target1"}).Return(nil)

	err := fallback.BuildTargets(context.Background(), internal.ProjectBuilderOptions{}, []string{"target1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestBuilderFallback_BuildTargets_Available2(t *testing.T) {
	tool1 := createBuilderTool("test1", nil)
	tool2 := createBuilderTool("test2", &internal.Executable{Path: "test1"})

	fallback := &BuilderFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("BuildTargets", mock.Anything, mock.Anything, []string{"target1"}).Return(nil)

	err := fallback.BuildTargets(context.Background(), internal.ProjectBuilderOptions{}, []string{"target1"})

	assert.NoError(t, err)

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

func TestTesterFallback_TestAll_Empty(t *testing.T) {
	fallback := &TesterFallback{}

	err := fallback.TestAll(context.Background(), internal.ProjectTesterOptions{})

	assert.EqualError(t, err, "no tester tool available: none")
}

func TestTesterFallback_TestAll_NoAvailable(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", nil)

	fallback := &TesterFallback{tool1, tool2}

	err := fallback.TestAll(context.Background(), internal.ProjectTesterOptions{})

	assert.EqualError(t, err, "no tester tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_TestAll_Available1(t *testing.T) {
	tool1 := createTesterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createTesterTool("test2", nil)

	fallback := &TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("TestAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.TestAll(context.Background(), internal.ProjectTesterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_TestAll_Available2(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("TestAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.TestAll(context.Background(), internal.ProjectTesterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_Test_Empty(t *testing.T) {
	fallback := &TesterFallback{}

	err := fallback.TestPattern(context.Background(), internal.ProjectTesterOptions{}, "target1")

	assert.EqualError(t, err, "no tester tool available: none")
}

func TestTesterFallback_Test_NoAvailable(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", nil)

	fallback := &TesterFallback{tool1, tool2}

	err := fallback.TestPattern(context.Background(), internal.ProjectTesterOptions{}, "target1")

	assert.EqualError(t, err, "no tester tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_Test_Available1(t *testing.T) {
	tool1 := createTesterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createTesterTool("test2", nil)

	fallback := &TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("TestPattern", mock.Anything, mock.Anything, "target1").Return(nil)

	err := fallback.TestPattern(context.Background(), internal.ProjectTesterOptions{}, "target1")

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestTesterFallback_Test_Available2(t *testing.T) {
	tool1 := createTesterTool("test1", nil)
	tool2 := createTesterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &TesterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("TestPattern", mock.Anything, mock.Anything, "target1").Return(nil)

	err := fallback.TestPattern(context.Background(), internal.ProjectTesterOptions{}, "target1")

	assert.NoError(t, err)

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

func TestFormatterFallback_FormatAll_Empty(t *testing.T) {
	fallback := &FormatterFallback{}

	err := fallback.FormatAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatAll_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	err := fallback.FormatAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatAll_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatAll_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatFiles_Empty(t *testing.T) {
	fallback := &FormatterFallback{}

	err := fallback.FormatFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatFiles_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	err := fallback.FormatFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatFiles_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatFiles_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckAll_Empty(t *testing.T) {
	fallback := &FormatterFallback{}

	err := fallback.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatCheckAll_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	err := fallback.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckAll_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatCheckAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckAll_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatCheckAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckFiles_Empty(t *testing.T) {
	fallback := &FormatterFallback{}

	err := fallback.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no formatter tool available: none")
}

func TestFormatterFallback_FormatCheckFiles_NoAvailable(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	err := fallback.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no formatter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckFiles_Available1(t *testing.T) {
	tool1 := createFormatterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createFormatterTool("test2", nil)

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FormatCheckFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestFormatterFallback_FormatCheckFiles_Available2(t *testing.T) {
	tool1 := createFormatterTool("test1", nil)
	tool2 := createFormatterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &FormatterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FormatCheckFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{}, []string{"file1"})

	assert.NoError(t, err)

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

func TestLinterFallback_LintAll_Empty(t *testing.T) {
	fallback := &LinterFallback{}

	err := fallback.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.EqualError(t, err, "no linter tool available: none")
}

func TestLinterFallback_LintAll_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	fallback := &LinterFallback{tool1, tool2}

	err := fallback.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintAll_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	fallback := &LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintAll_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := fallback.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintFiles_Empty(t *testing.T) {
	fallback := &LinterFallback{}

	err := fallback.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no linter tool available: none")
}

func TestLinterFallback_LintFiles_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	fallback := &LinterFallback{tool1, tool2}

	err := fallback.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintFiles_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	fallback := &LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterFallback_LintFiles_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	fallback := &LinterFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := fallback.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintAll_Empty(t *testing.T) {
	list := &LinterList{}

	err := list.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.EqualError(t, err, "no linter tool available: none")
}

func TestLinterList_LintAll_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	list := &LinterList{tool1, tool2}

	err := list.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintAll_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	list := &LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := list.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintAll_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	list := &LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintAll", mock.Anything, mock.Anything).Return(nil)
	tool2.On("LintAll", mock.Anything, mock.Anything).Return(nil)

	err := list.LintAll(context.Background(), internal.ProjectLinterOptions{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintFiles_Empty(t *testing.T) {
	list := &LinterList{}

	err := list.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no linter tool available: none")
}

func TestLinterList_LintFiles_NoAvailable(t *testing.T) {
	tool1 := createLinterTool("test1", nil)
	tool2 := createLinterTool("test2", nil)

	list := &LinterList{tool1, tool2}

	err := list.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.EqualError(t, err, "no linter tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintFiles_Available1(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", nil)

	list := &LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := list.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestLinterList_LintFiles_Available2(t *testing.T) {
	tool1 := createLinterTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createLinterTool("test2", &internal.Executable{Path: "test1"})

	list := &LinterList{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)
	tool2.On("LintFiles", mock.Anything, mock.Anything, []string{"file1"}).Return(nil)

	err := list.LintFiles(context.Background(), internal.ProjectLinterOptions{}, []string{"file1"})

	assert.NoError(t, err)

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

func TestRunnerFallback_RunTarget_Empty(t *testing.T) {
	fallback := &RunnerFallback{}

	err := fallback.RunTarget(context.Background(), internal.ProjectInfo{}, "target1", []string{})

	assert.EqualError(t, err, "no runner tool available: none")
}

func TestRunnerFallback_RunTarget_NoAvailable(t *testing.T) {
	tool1 := createRunnerTool("test1", nil)
	tool2 := createRunnerTool("test2", nil)

	fallback := &RunnerFallback{tool1, tool2}

	err := fallback.RunTarget(context.Background(), internal.ProjectInfo{}, "target1", []string{})

	assert.EqualError(t, err, "no runner tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestRunnerFallback_RunTarget_Available1(t *testing.T) {
	tool1 := createRunnerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createRunnerTool("test2", nil)

	fallback := &RunnerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("RunTarget", mock.Anything, mock.Anything, "target1", []string{}).Return(nil)

	err := fallback.RunTarget(context.Background(), internal.ProjectInfo{}, "target1", []string{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestRunnerFallback_RunTarget_Available2(t *testing.T) {
	tool1 := createRunnerTool("test1", nil)
	tool2 := createRunnerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &RunnerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("RunTarget", mock.Anything, mock.Anything, "target1", []string{}).Return(nil)

	err := fallback.RunTarget(context.Background(), internal.ProjectInfo{}, "target1", []string{})

	assert.NoError(t, err)

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

func TestDependencyManagerFallback_AddDependencies_Empty(t *testing.T) {
	fallback := &DependencyManagerFallback{}

	err := fallback.AddDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_AddDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	err := fallback.AddDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AddDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("AddDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.AddDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AddDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("AddDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.AddDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_RemoveDependencies_Empty(t *testing.T) {
	fallback := &DependencyManagerFallback{}

	err := fallback.RemoveDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_RemoveDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	err := fallback.RemoveDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_RemoveDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("RemoveDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.RemoveDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_RemoveDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("RemoveDependencies", mock.Anything, mock.Anything, []string{"dep1"}, false).Return(nil)

	err := fallback.RemoveDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"}, false)

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_UpdateDependencies_Empty(t *testing.T) {
	fallback := &DependencyManagerFallback{}

	err := fallback.UpdateDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"})

	assert.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_UpdateDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	err := fallback.UpdateDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"})

	assert.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_UpdateDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("UpdateDependencies", mock.Anything, mock.Anything, []string{"dep1"}).Return(nil)

	err := fallback.UpdateDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_UpdateDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("UpdateDependencies", mock.Anything, mock.Anything, []string{"dep1"}).Return(nil)

	err := fallback.UpdateDependencies(context.Background(), internal.ProjectInfo{}, []string{"dep1"})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_FetchDependencies_Empty(t *testing.T) {
	fallback := &DependencyManagerFallback{}

	err := fallback.FetchDependencies(context.Background(), internal.ProjectInfo{}, false)

	assert.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_FetchDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	err := fallback.FetchDependencies(context.Background(), internal.ProjectInfo{}, false)

	assert.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_FetchDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("FetchDependencies", mock.Anything, mock.Anything, false).Return(nil)

	err := fallback.FetchDependencies(context.Background(), internal.ProjectInfo{}, false)

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_FetchDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("FetchDependencies", mock.Anything, mock.Anything, false).Return(nil)

	err := fallback.FetchDependencies(context.Background(), internal.ProjectInfo{}, false)

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_ListDependencies_Empty(t *testing.T) {
	fallback := &DependencyManagerFallback{}

	err := fallback.ListDependencies(context.Background(), internal.ProjectInfo{})

	assert.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_ListDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	err := fallback.ListDependencies(context.Background(), internal.ProjectInfo{})

	assert.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_ListDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("ListDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.ListDependencies(context.Background(), internal.ProjectInfo{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_ListDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("ListDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.ListDependencies(context.Background(), internal.ProjectInfo{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AuditDependencies_Empty(t *testing.T) {
	fallback := &DependencyManagerFallback{}

	err := fallback.AuditDependencies(context.Background(), internal.ProjectInfo{})

	assert.EqualError(t, err, "no dependency management tool available: none")
}

func TestDependencyManagerFallback_AuditDependencies_NoAvailable(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	err := fallback.AuditDependencies(context.Background(), internal.ProjectInfo{})

	assert.EqualError(t, err, "no dependency management tool available: test1, test2")

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AuditDependencies_Available1(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", &internal.Executable{Path: "test1"})
	tool2 := createDependencyManagerTool("test2", nil)

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool1.On("AuditDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.AuditDependencies(context.Background(), internal.ProjectInfo{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}

func TestDependencyManagerFallback_AuditDependencies_Available2(t *testing.T) {
	tool1 := createDependencyManagerTool("test1", nil)
	tool2 := createDependencyManagerTool("test2", &internal.Executable{Path: "test1"})

	fallback := &DependencyManagerFallback{tool1, tool2}

	tool1.Test(t)
	tool2.Test(t)
	tool2.On("AuditDependencies", mock.Anything, mock.Anything).Return(nil)

	err := fallback.AuditDependencies(context.Background(), internal.ProjectInfo{})

	assert.NoError(t, err)

	tool1.AssertExpectations(t)
	tool2.AssertExpectations(t)
}
