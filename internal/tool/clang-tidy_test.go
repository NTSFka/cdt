package tool_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClangTidy_DetectClangTidy(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-tidy").
		Return(env.NewExecutable("/bin/clang-tidy"), nil)

	clangTidy := tool.DetectClangTidy(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, clangTidy)
	assert.Equal(t, "clang-tidy", clangTidy.Id())
	assert.True(t, clangTidy.IsAvailable())

	if executable := clangTidy.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/clang-tidy", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestClangTidy_DetectClangTidy_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-tidy").
		Return(nil, nil)

	clangTidy := tool.DetectClangTidy(t.Context(), tool.DetectOptions{Environment: env})
	assert.NotNil(t, clangTidy)
	assert.Equal(t, "clang-tidy", clangTidy.Id())
	assert.False(t, clangTidy.IsAvailable())
	assert.Nil(t, clangTidy.Executable())

	env.AssertExpectations(t)
}

func TestClangTidy_DetectClangTidy_Config(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-tidy-20").
		Return(env.NewExecutable("/bin/clang-tidy"), nil)

	clangTidy := tool.DetectClangTidy(t.Context(), tool.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"clang-tidy": "clang-tidy-20"},
	})
	assert.NotNil(t, clangTidy)
	assert.Equal(t, "clang-tidy", clangTidy.Id())
	assert.True(t, clangTidy.IsAvailable())

	if executable := clangTidy.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/clang-tidy", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestClangTidy_LintFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:       "project",
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.c", "file2.c"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *info.OutputDirectory,
		"file1.c",
		"file2.c",
	}).
		Return(nil)

	err := clangTidy.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_All_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:       "project",
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.c", "file2.c"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *info.OutputDirectory,
		"file1.c",
		"file2.c",
	}).
		Return(errors.New("failed"))

	err := clangTidy.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_All_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.c", "file2.c"},
					},
				},
			},
		},
	}

	file, err := os.Create(filepath.Join(info.Directory, ".clang-tidy"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(info.Directory, ".clang-tidy")),
		"-p", *info.OutputDirectory,
		"file1.c",
		"file2.c",
	}).
		Return(nil)

	err = clangTidy.LintFiles(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *info.OutputDirectory,
		"file1.c",
		filepath.Join(info.Directory, "file3.c"),
	}).
		Return(nil)

	err := clangTidy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file1.c", filepath.Join(info.Directory, "file3.c")},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:       "project",
		OutputDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *info.OutputDirectory,
		"file1.c",
	}).
		Return(errors.New("failed"))

	err := clangTidy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info, Filenames: &[]string{"file1.c"}},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
	}

	file, err := os.Create(filepath.Join(info.Directory, ".clang-tidy"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(info.Directory, ".clang-tidy")),
		"-p", *info.OutputDirectory,
		"file1.c",
	}).
		Return(nil)

	err = clangTidy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info, Filenames: &[]string{"file1.c"}},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_ClangTidy_LintFiles_OutputFormat_Raw(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.c", "file2.c"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *info.OutputDirectory,
		"file1.c",
		"file2.c",
	}).
		Return(nil)

	err := clangTidy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Format: internal.LintReportFormatRaw,
			},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_ClangTidy_LintFiles_OutputFormat_Unsupported(t *testing.T) {
	dataSet := []struct {
		Format internal.LintReportFormat
	}{
		{internal.LintReportFormatJson},
		{internal.LintReportFormatJUnit},
		{internal.LintReportFormatGitHub},
		{internal.LintReportFormatGitLab},
		{internal.LintReportFormatTeamCity},
		{"test-unsupported"},
	}

	for _, data := range dataSet {
		t.Run(string(data.Format), func(t *testing.T) {
			exec := test.NewExecutable(t)

			clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

			info := internal.ProjectInfo{
				Directory:       t.TempDir(),
				OutputDirectory: internal.StrPtr("build"),
				StructureProvider: &internal.FixedProjectStructureProvider{
					ProjectStructure: internal.ProjectStructure{
						Targets: map[string]internal.ProjectTarget{
							"target1": {
								Files: []string{"file1.c", "file2.c"},
							},
						},
					},
				},
			}

			err := clangTidy.LintFiles(
				t.Context(),
				internal.ProjectLinterOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: data.Format,
					},
				},
			)
			require.EqualError(t, err, "unsupported report format: "+string(data.Format))

			exec.AssertExpectations(t)
		})
	}
}

func TestClangTidy_ClangTidy_LintFiles_OutputFile(t *testing.T) {
	exec := test.NewExecutable(t)

	linter := tool.NewClangTidy(exec.LazyExecutable("lint"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.c", "file2.c"},
					},
				},
			},
		},
	}

	exec.OnRunOutput("lint", []string{
		"-p", *info.OutputDirectory,
		"file1.c",
		"file2.c",
	}, "ok").
		Return(nil)

	outputFilename := t.TempDir() + "/report.txt"

	err := linter.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Filename: &outputFilename,
			},
		},
	)
	require.NoError(t, err)
	require.FileExists(t, outputFilename)

	outputContent, err := os.ReadFile(outputFilename)
	require.NoError(t, err)
	assert.Equal(t, "ok", string(outputContent))

	exec.AssertExpectations(t)
}

func TestClangTidy_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.ProjectInfo{Directory: "build"}

	exec.OnRun("clang-tidy", []string{desc.Directory}).
		Return(nil)

	err := clangTidy.RunForProject(t.Context(), desc, []string{})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	desc := internal.ProjectInfo{Directory: "build"}

	exec.OnRun("clang-tidy", []string{desc.Directory}).
		Return(errors.New("failed"))

	err := clangTidy.RunForProject(t.Context(), desc, []string{})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
