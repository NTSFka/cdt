package tool_test

import (
	"context"
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
	goTool := test.RunDetectToolLastFound(t, goDetect, []string{"go"})

	assert.Equal(t, tool.IdGo, goTool.Id())
}

func TestGo_DetectGo_NotFound(t *testing.T) {
	goTool := test.RunDetectToolNotFound(t, goDetect, []string{"go"})

	assert.Equal(t, tool.IdGo, goTool.Id())
}

func TestGo_DetectGo_Config(t *testing.T) {
	goTool := test.RunDetectToolConfig(t, goDetect, "go")

	assert.Equal(t, tool.IdGo, goTool.Id())
}

func TestGo_Structure(t *testing.T) {
	test.InvokeStructure(
		t,
		goStructureProvider,
		internal.ProjectInfo{Directory: "."},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(
				executablePath,
				[]string{"list", "-json=ImportPath,GoFiles", "./..."},
				`
				{"ImportPath": "target1","GoFiles":["file1.go"]}
				{"ImportPath": "target2","GoFiles":["file2.go", "file3.go"]}
				`,
			).
				Return(nil)
		},
		func(structure *internal.ProjectStructure, err error) {
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
		},
	)
}

func TestGo_Structure_Failed(t *testing.T) {
	test.InvokeStructure(
		t,
		goStructureProvider,
		internal.ProjectInfo{Directory: "."},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{"list", "-json=ImportPath,GoFiles", "./..."}).
				Return(errors.New("failed"))
		},
		func(structure *internal.ProjectStructure, err error) {
			require.EqualError(t, err, "failed")
			assert.Nil(t, structure)
		},
	)
}

func TestGo_Structure_InvalidJson(t *testing.T) {
	test.InvokeStructure(
		t,
		goStructureProvider,
		internal.ProjectInfo{Directory: "."},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"list", "-json=ImportPath,GoFiles", "./..."}, `{]}`).
				Return(nil)
		},
		func(structure *internal.ProjectStructure, err error) {
			require.ErrorContains(t, err, "json decode failed:")
			assert.Nil(t, structure)
		},
	)
}

func TestGo_BuildTargets_All(t *testing.T) {
	test.InvokeBuildTargetsSuccess(
		t,
		goBuilder,
		internal.ProjectBuilderOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"build"},
	)
}

func TestGo_BuildTargets_All_Failed(t *testing.T) {
	test.InvokeBuildTargetsFailed(
		t,
		goBuilder,
		internal.ProjectBuilderOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"build"},
	)
}

func TestGo_BuildTargets_All_OutputDirectory(t *testing.T) {
	outputDir := t.TempDir()

	test.InvokeBuildTargetsSuccess(
		t,
		goBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       ".",
				OutputDirectory: &outputDir,
			},
		},
		[]string{"build", "-o", outputDir},
	)
}

func TestGo_BuildTargets(t *testing.T) {
	test.InvokeBuildTargetsSuccess(
		t,
		goBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Targets:     &[]string{"target1", "target2"},
		},
		[]string{"build", "target1", "target2"},
	)
}

func TestGo_BuildTargets_Failed(t *testing.T) {
	test.InvokeBuildTargetsFailed(
		t,
		goBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Targets:     &[]string{"target1", "target2"},
		},
		[]string{"build", "target1", "target2"},
	)
}

func TestGo_BuildTargets_OutputDirectory(t *testing.T) {
	outputDir := t.TempDir()

	test.InvokeBuildTargetsSuccess(
		t,
		goBuilder,
		internal.ProjectBuilderOptions{
			ProjectInfo: internal.ProjectInfo{Directory: ".", OutputDirectory: &outputDir},
			Targets:     &[]string{"target1", "target2"},
		},
		[]string{"build", "-o", outputDir, "target1", "target2"},
	)
}

func TestGo_RunTarget(t *testing.T) {
	test.InvokeRunTargetSuccess(
		t,
		goRunner,
		internal.ProjectRunnerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		"target1",
		[]string{"run", "target1"},
	)
}

func TestGo_RunTarget_Failed(t *testing.T) {
	test.InvokeRunTargetFailed(
		t,
		goRunner,
		internal.ProjectRunnerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		"target1",
		[]string{"run", "target1"},
	)
}

func TestGo_RunTests(t *testing.T) {
	test.InvokeRunTestsSuccess(
		t,
		goTester,
		internal.ProjectTesterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"test", "./..."},
	)
}

func TestGo_RunTests_Failed(t *testing.T) {
	test.InvokeRunTestsFailed(
		t,
		goTester,
		internal.ProjectTesterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"test", "./..."},
	)
}

func TestGo_RunTests_Pattern(t *testing.T) {
	test.InvokeRunTestsSuccess(
		t,
		goTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     internal.StrPtr("test1"),
		},
		[]string{"test", "test1"},
	)
}

func TestGo_RunTests_Pattern_Failed(t *testing.T) {
	test.InvokeRunTestsFailed(
		t,
		goTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Pattern:     internal.StrPtr("test1"),
		},
		[]string{"test", "test1"},
	)
}

func TestGo_RunTests_Pattern_FormatRawEvents(t *testing.T) {
	tempDir := t.TempDir()
	reportFilename := filepath.Join(tempDir, "report.json")

	test.InvokeRunTests(
		t,
		goTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatRawEvents,
				Filename: &reportFilename,
			},
			Pattern: internal.StrPtr("test1"),
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"test", "-json", "test1"}, `
				{"Time":"2026-01-25T15:57:17.480674482+01:00","Action":"start","Package":"t"}
				{"Time":"2026-01-25T15:57:17.480691796+01:00","Action":"run","Package":"t","Test":"Test1"}
				{"Time":"2026-01-25T15:57:17.480717611+01:00","Action":"pass","Package":"t","Test":"Test1","Elapsed":0}
			`).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)

			report, err := os.ReadFile(reportFilename)
			require.NoError(t, err)
			assert.NotEmpty(t, string(report))
		},
	)
}

func TestGo_RunTests_Pattern_FormatCtrf(t *testing.T) {
	tempDir := t.TempDir()
	reportFilename := filepath.Join(tempDir, "report.json")

	test.InvokeRunTests(
		t,
		goTester,
		internal.ProjectTesterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Output: internal.OutputOptions[internal.TestsReportFormat]{
				Format:   internal.TestsReportFormatCtrf,
				Filename: &reportFilename,
			},
			Pattern: internal.StrPtr("test1"),
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{"test", "-json", "test1"}, `
				{"Time":"2026-01-25T15:57:17.480674482+01:00","Action":"start","Package":"t"}
				{"Time":"2026-01-25T15:57:17.480691796+01:00","Action":"run","Package":"t","Test":"Test1"}
				{"Time":"2026-01-25T15:57:17.480717611+01:00","Action":"pass","Package":"t","Test":"Test1","Elapsed":0}
			`).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)

			report, err := os.ReadFile(reportFilename)
			require.NoError(t, err)

			assert.Contains(t, string(report), `"Test1"`)
		},
	)
}

func TestGo_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		goFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"fmt", "./..."},
	)
}

func TestGo_FormatFiles_All_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		goFormatter,
		internal.ProjectFormatterOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"fmt", "./..."},
	)
}

func TestGo_FormatFiles(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		goFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Filenames:   &[]string{"file1"},
		},
		[]string{"fmt", "file1"},
	)
}

func TestGo_FormatFiles_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		goFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			Filenames:   &[]string{"file1"},
		},
		[]string{"fmt", "file1"},
	)
}

func TestGo_FormatFiles_CheckAll(t *testing.T) {
	test.InvokeFormatFiles(
		t,
		goFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
		},
		func(exec *test.Executable, executablePath string) {
		},
		func(err error) {
			require.EqualError(t, err, "go fmt doesn't support check mode")
		},
	)
}

func TestGo_FormatFilesFiles_Check(t *testing.T) {
	test.InvokeFormatFiles(
		t,
		goFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
			Filenames:   &[]string{"file1"},
		},
		func(exec *test.Executable, executablePath string) {
		},
		func(err error) {
			require.EqualError(t, err, "go fmt doesn't support check mode")
		},
	)
}

func TestGo_Go_LintFiles_All(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		goLinter,
		internal.ProjectLinterOptions{},
		[]string{"vet", "./..."},
	)
}

func TestGo_Go_LintFiles_All_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		goLinter,
		internal.ProjectLinterOptions{},
		[]string{"vet", "./..."},
	)
}

func TestGo_Go_LintFiles(t *testing.T) {
	test.InvokeLintFilesSuccess(
		t,
		goLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"mod1"}},
		[]string{"vet", "mod1"},
	)
}

func TestGo_Go_LintFiles_Failed(t *testing.T) {
	test.InvokeLintFilesFailed(
		t,
		goLinter,
		internal.ProjectLinterOptions{Filenames: &[]string{"mod1"}},
		[]string{"vet", "mod1"},
	)
}

func TestGo_Go_LintFiles_OutputFormat_Raw(t *testing.T) {
	test.InvokeLintFilesOutputFormatCheck(
		t,
		goLinter,
		".",
		internal.LintReportFormatRaw,
		[]string{"vet", "./..."},
		nil,
	)
}

func TestGo_Go_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []internal.LintReportFormat{
		internal.LintReportFormatJson,
		internal.LintReportFormatJUnit,
		internal.LintReportFormatGitHub,
		internal.LintReportFormatGitLab,
		internal.LintReportFormatTeamCity,
		"test-unsupported",
	}

	for _, format := range dataSet {
		t.Run(string(format), func(t *testing.T) {
			test.InvokeLintFilesOutputFormatUnsupported(t, goLinter, ".", format, nil)
		})
	}
}

func TestGo_Go_LintFiles_OutputFile(t *testing.T) {
	test.InvokeLintFilesOutputFile(t, goLinter, ".", []string{"vet", "./..."}, nil)
}

func TestGo_Go_AddDependencies(t *testing.T) {
	test.InvokeAddDependenciesSuccess(
		t,
		goDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		false,
		[]string{"get", "dep1"},
	)
}

func TestGo_Go_RemoveDependencies(t *testing.T) {
	test.InvokeRemoveDependenciesSuccess(
		t,
		goDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		false,
		[]string{"get", "dep1@none"},
	)
}

func TestGo_Go_UpdateDependencies(t *testing.T) {
	test.InvokeUpdateDependenciesSuccess(
		t,
		goDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"dep1"},
		[]string{"get", "dep1"},
	)
}

func TestGo_Go_FetchDependencies(t *testing.T) {
	test.InvokeFetchDependenciesSuccess(
		t,
		goDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		false,
		[]string{"mod", "tidy"},
	)
}

func TestGo_Go_ListDependencies(t *testing.T) {
	test.InvokeListDependenciesSuccess(
		t,
		goDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		[]string{"list", "-m", "all"},
	)
}

func TestGo_Go_AuditDependencies(t *testing.T) {
	test.InvokeAuditDependencies(
		t,
		goDependencyManager,
		internal.ProjectDependencyManagerOptions{ProjectInfo: internal.ProjectInfo{Directory: "."}},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{"audit"}).
				Return(nil)
		},
		func(err error) {
			require.EqualError(t, err, "not supported")
		},
	)
}

func goDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectGo(ctx, options)
}

func goStructureProvider(
	executable func() (*internal.Executable, error),
) internal.ProjectStructureProvider {
	return tool.NewGo(executable)
}

func goLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewGo(executable)
}

func goFormatter(executable func() (*internal.Executable, error)) internal.ProjectFormatter {
	return tool.NewGo(executable)
}

func goBuilder(executable func() (*internal.Executable, error)) internal.ProjectBuilder {
	return tool.NewGo(executable)
}

func goRunner(executable func() (*internal.Executable, error)) internal.ProjectRunner {
	return tool.NewGo(executable)
}

func goTester(executable func() (*internal.Executable, error)) internal.ProjectTester {
	return tool.NewGo(executable)
}

func goDependencyManager(
	executable func() (*internal.Executable, error),
) internal.ProjectDependencyManager {
	return tool.NewGo(executable)
}
