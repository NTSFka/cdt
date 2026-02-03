package tool_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGo_DetectGo(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("go").
		Return(env.NewExecutable("/bin/go"), nil)

	goTool := tool.DetectGo(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, goTool)
	assert.Equal(t, "go", goTool.Id())
	assert.True(t, goTool.IsAvailable())

	if executable := goTool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/go", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestGo_DetectGo_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("go").
		Return(nil, nil)

	goTool := tool.DetectGo(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, goTool)
	assert.Equal(t, "go", goTool.Id())
	assert.False(t, goTool.IsAvailable())
	assert.Nil(t, goTool.Executable())

	env.AssertExpectations(t)
}

func TestGo_DetectGo_Config(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("go-1").
		Return(env.NewExecutable("/bin/go"), nil)

	goTool := tool.DetectGo(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"go": "go-1"},
	})
	assert.NotNil(t, goTool)
	assert.Equal(t, "go", goTool.Id())
	assert.True(t, goTool.IsAvailable())

	if executable := goTool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/go", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestGo_Structure(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRunOutput(
		"go",
		[]string{"list", "-json=ImportPath,GoFiles", "./..."},
		`{"ImportPath": "target1","GoFiles":["file1.go"]}{"ImportPath": "target2","GoFiles":["file2.go", "file3.go"]}`,
	).
		Return(nil)

	structure, err := goTool.Structure(t.Context(), info)
	require.NoError(t, err)
	require.NotNil(t, structure)
	assert.Equal(t,
		internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				"target1": {
					Files: []string{"file1.go"},
				},
				"target2": {
					Files: []string{"file2.go", "file3.go"},
				},
			},
		},
		*structure,
	)

	exec.AssertExpectations(t)
}

func TestGo_Structure_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"list", "-json=ImportPath,GoFiles", "./..."}).
		Return(errors.New("failed"))

	structure, err := goTool.Structure(t.Context(), info)
	require.EqualError(t, err, "failed")
	assert.Nil(t, structure)

	exec.AssertExpectations(t)
}

func TestGo_Structure_InvalidJson(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRunOutput("go", []string{"list", "-json=ImportPath,GoFiles", "./..."}, `{]}`).
		Return(nil)

	structure, err := goTool.Structure(t.Context(), info)
	require.ErrorContains(t, err, "json decode failed:")
	assert.Nil(t, structure)

	exec.AssertExpectations(t)
}

func TestGo_BuildAll(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build"}).
		Return(nil)

	err := goTool.BuildAll(t.Context(), internal.ProjectBuilderOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_BuildAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build"}).
		Return(errors.New("failed"))

	err := goTool.BuildAll(t.Context(), internal.ProjectBuilderOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_BuildAll_OutputDirectory(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	outputDir := t.TempDir()

	info := internal.ProjectInfo{Directory: ".", OutputDirectory: &outputDir}

	exec.OnRun("go", []string{"build", "-o", outputDir}).
		Return(nil)

	err := goTool.BuildAll(t.Context(), internal.ProjectBuilderOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_BuildTargets(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build", "target1", "target2"}).
		Return(nil)

	err := goTool.BuildTargets(
		t.Context(),
		internal.ProjectBuilderOptions{ProjectInfo: info},
		[]string{"target1", "target2"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_BuildTargets_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"build", "target1", "target2"}).
		Return(errors.New("failed"))

	err := goTool.BuildTargets(
		t.Context(),
		internal.ProjectBuilderOptions{ProjectInfo: info},
		[]string{"target1", "target2"},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_BuildTargets_OutputDirectory(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	outputDir := t.TempDir()

	info := internal.ProjectInfo{Directory: ".", OutputDirectory: &outputDir}

	exec.OnRun("go", []string{"build", "-o", outputDir, "target1", "target2"}).
		Return(nil)

	err := goTool.BuildTargets(
		t.Context(),
		internal.ProjectBuilderOptions{ProjectInfo: info},
		[]string{"target1", "target2"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_RunTarget(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"run", "target1"}).
		Return(nil)

	err := goTool.RunTarget(
		t.Context(),
		internal.ProjectRunnerOptions{ProjectInfo: info},
		"target1",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_RunTarget_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"run", "target1"}).
		Return(errors.New("failed"))

	err := goTool.RunTarget(
		t.Context(),
		internal.ProjectRunnerOptions{ProjectInfo: info},
		"target1",
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_TestAll(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "./..."}).
		Return(nil)

	err := goTool.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_TestAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "./..."}).
		Return(errors.New("failed"))

	err := goTool.TestAll(t.Context(), internal.ProjectTesterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_Test(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "test1"}).
		Return(nil)

	err := goTool.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: info},
		"test1",
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Test_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"test", "test1"}).
		Return(errors.New("failed"))

	err := goTool.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{ProjectInfo: info},
		"test1",
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_TestPattern_FormatUnsupported(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	data := []internal.TestsReportFormat{
		"test",
		internal.TestsReportFormatJson,
		internal.TestsReportFormatJUnit,
		internal.TestsReportFormatTeamCity,
	}

	for _, format := range data {
		t.Run(string(format), func(t *testing.T) {
			err := goTool.TestPattern(
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

func TestGo_TestPattern_FormatRawEvents(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRunOutput("go", []string{"test", "-json", "test1"}, `
		{"Time":"2026-01-25T15:57:17.480674482+01:00","Action":"start","Package":"t"}
		{"Time":"2026-01-25T15:57:17.480691796+01:00","Action":"run","Package":"t","Test":"Test1"}
		{"Time":"2026-01-25T15:57:17.480717611+01:00","Action":"pass","Package":"t","Test":"Test1","Elapsed":0}
	`).
		Return(nil)

	tempDir := t.TempDir()
	reportFilename := filepath.Join(tempDir, "report.json")

	err := goTool.TestPattern(
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

	report, err := os.ReadFile(reportFilename)
	require.NoError(t, err)
	assert.NotEmpty(t, string(report))

	exec.AssertExpectations(t)
}

func TestGo_TestPattern_FormatCtrf(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRunOutput("go", []string{"test", "-json", "test1"}, `
		{"Time":"2026-01-25T15:57:17.480674482+01:00","Action":"start","Package":"t"}
		{"Time":"2026-01-25T15:57:17.480691796+01:00","Action":"run","Package":"t","Test":"Test1"}
		{"Time":"2026-01-25T15:57:17.480717611+01:00","Action":"pass","Package":"t","Test":"Test1","Elapsed":0}
	`).
		Return(nil)

	tempDir := t.TempDir()
	reportFilename := filepath.Join(tempDir, "report.json")

	err := goTool.TestPattern(
		t.Context(),
		internal.ProjectTesterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatCtrf,
				Filename: &reportFilename,
			},
		},
		"test1",
	)
	require.NoError(t, err)

	report, err := os.ReadFile(reportFilename)
	require.NoError(t, err)

	assert.Contains(t, string(report), `"Test1"`)

	exec.AssertExpectations(t)
}

func TestGo_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "./..."}).
		Return(nil)

	err := goTool.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_FormatAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "./..."}).
		Return(errors.New("failed"))

	err := goTool.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "file1"}).
		Return(nil)

	err := goTool.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"file1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_FormatFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"fmt", "file1"}).
		Return(errors.New("failed"))

	err := goTool.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"file1"},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestGo_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	err := goTool.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.EqualError(t, err, "go fmt doesn't support check mode")

	exec.AssertExpectations(t)
}

func TestGo_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	err := goTool.FormatCheckFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info},
		[]string{"file1"},
	)
	require.EqualError(t, err, "go fmt doesn't support check mode")

	exec.AssertExpectations(t)
}

func TestGo_Go_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"vet", "./..."}).
		Return(nil)

	err := goTool.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_Lint(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("lint", []string{"vet", "mod1"}).
		Return(nil)

	err := goTool.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"mod1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_AddDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"get", "dep1"}).
		Return(nil)

	err := goTool.AddDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_RemoveDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"get", "dep1@none"}).
		Return(nil)

	err := goTool.RemoveDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_UpdateDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"get", "dep1"}).
		Return(nil)

	err := goTool.UpdateDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		[]string{"dep1"},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_FetchDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"mod", "tidy"}).
		Return(nil)

	err := goTool.FetchDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
		false,
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_ListDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"list", "-m", "all"}).
		Return(nil)

	err := goTool.ListDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestGo_Go_AuditDependencies(t *testing.T) {
	exec := test.NewExecutable(t)

	goTool := tool.NewGo(exec.LazyExecutable("go"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("go", []string{"audit"}).
		Return(nil)

	err := goTool.AuditDependencies(
		t.Context(),
		internal.ProjectDependencyManagerOptions{ProjectInfo: info},
	)
	require.EqualError(t, err, "not supported")

	exec.AssertExpectations(t)
}
