package tool_test

import (
	"context"
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
	clangTidy := test.RunDetectToolLastFound(t, clangTidyDetect, []string{"clang-tidy"})

	assert.Equal(t, tool.IdClangTidy, clangTidy.Id())
}

func TestClangTidy_DetectClangTidy_NotFound(t *testing.T) {
	clangTidy := test.RunDetectToolNotFound(t, clangTidyDetect, []string{"clang-tidy"})

	assert.Equal(t, tool.IdClangTidy, clangTidy.Id())
}

func TestClangTidy_DetectClangTidy_Config(t *testing.T) {
	clangTidy := test.RunDetectToolConfig(t, clangTidyDetect, "clang-tidy")

	assert.Equal(t, tool.IdClangTidy, clangTidy.Id())
}

func TestClangTidy_LintFiles_All(t *testing.T) {
	info := clangTidyProjectInfo(t, []string{"file1.c", "file2.c"})

	test.InvokeLintFiles(
		t,
		clangTidyLinter,
		internal.ProjectLinterOptions{ProjectInfo: info},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{
				"-p", *info.OutputDirectory,
				"file1.c",
				"file2.c",
			}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestClangTidy_LintFiles_All_Failed(t *testing.T) {
	info := clangTidyProjectInfo(t, []string{"file1.c", "file2.c"})

	test.InvokeLintFiles(
		t,
		clangTidyLinter,
		internal.ProjectLinterOptions{ProjectInfo: info},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{
				"-p", *info.OutputDirectory,
				"file1.c",
				"file2.c",
			}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func TestClangTidy_LintFiles_All_CustomConfig(t *testing.T) {
	info := clangTidyProjectInfo(t, []string{"file1.c", "file2.c"})

	// Create configuration file
	file, err := os.Create(filepath.Join(info.Directory, ".clang-tidy"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	test.InvokeLintFiles(
		t,
		clangTidyLinter,
		internal.ProjectLinterOptions{ProjectInfo: info},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{
				fmt.Sprintf("--config-file=%v", filepath.Join(info.Directory, ".clang-tidy")),
				"-p", *info.OutputDirectory,
				"file1.c",
				"file2.c",
			}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestClangTidy_LintFiles(t *testing.T) {
	info := clangTidyProjectInfo(t, []string{})

	test.InvokeLintFiles(
		t,
		clangTidyLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file1.c", filepath.Join(info.Directory, "file3.c")},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{
				"-p", *info.OutputDirectory,
				"file1.c",
				filepath.Join(info.Directory, "file3.c"),
			}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestClangTidy_LintFiles_Failed(t *testing.T) {
	info := clangTidyProjectInfo(t, []string{})

	test.InvokeLintFiles(
		t,
		clangTidyLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file1.c", filepath.Join(info.Directory, "file3.c")},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{
				"-p", *info.OutputDirectory,
				"file1.c",
				filepath.Join(info.Directory, "file3.c"),
			}).
				Return(errors.New("failed"))
		},
		func(err error) {
			require.EqualError(t, err, "failed")
		},
	)
}

func TestClangTidy_LintFiles_CustomConfig(t *testing.T) {
	info := clangTidyProjectInfo(t, []string{"file1.c", "file2.c"})

	// Create configuration file
	file, err := os.Create(filepath.Join(info.Directory, ".clang-tidy"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	test.InvokeLintFiles(
		t,
		clangTidyLinter,
		internal.ProjectLinterOptions{ProjectInfo: info, Filenames: &[]string{"file1.c"}},
		func(exec *test.Executable, executablePath string) {
			exec.OnRun(executablePath, []string{
				fmt.Sprintf("--config-file=%v", filepath.Join(info.Directory, ".clang-tidy")),
				"-p", *info.OutputDirectory,
				"file1.c",
			}).
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
		},
	)
}

func TestClangTidy_ClangTidy_LintFiles_OutputFormat(t *testing.T) {
	dataSet := []struct {
		format internal.LintReportFormat
		args   []string
	}{
		{internal.LintReportFormatRaw, []string{}},
	}

	for _, data := range dataSet {
		t.Run(string(data.format), func(t *testing.T) {
			info := clangTidyProjectInfo(t, []string{"file1.c", "file2.c"})

			test.InvokeLintFiles(
				t,
				clangTidyLinter,
				internal.ProjectLinterOptions{
					ProjectInfo: info,
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: data.format,
					},
				},
				func(exec *test.Executable, executablePath string) {
					exec.OnRun(executablePath, append([]string{
						"-p", *info.OutputDirectory,
						"file1.c",
						"file2.c",
					}, data.args...)).
						Return(nil)
				},
				func(err error) {
					require.NoError(t, err)
				},
			)
		})
	}
}

func TestClangTidy_ClangTidy_LintFiles_OutputFormat_Unsupported(t *testing.T) {
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
			test.InvokeLintFiles(
				t,
				clangTidyLinter,
				internal.ProjectLinterOptions{
					ProjectInfo: clangTidyProjectInfo(t, []string{"file1.c", "file2.c"}),
					Output: internal.OutputOptions[internal.LintReportFormat]{
						Format: format,
					},
				},
				func(exec *test.Executable, executablePath string) {
					// Not executed
				},
				func(err error) {
					require.EqualError(t, err, "unsupported report format: "+string(format))
				},
			)
		})
	}
}

func TestClangTidy_ClangTidy_LintFiles_OutputFile(t *testing.T) {
	info := clangTidyProjectInfo(t, []string{"file1.c", "file2.c"})

	outputFilename := t.TempDir() + "/report.txt"

	test.InvokeLintFiles(
		t,
		clangTidyLinter,
		internal.ProjectLinterOptions{
			ProjectInfo: info,
			Output: internal.OutputOptions[internal.LintReportFormat]{
				Filename: &outputFilename,
			},
		},
		func(exec *test.Executable, executablePath string) {
			exec.OnRunOutput(executablePath, []string{
				"-p", *info.OutputDirectory,
				"file1.c",
				"file2.c",
			}, "ok").
				Return(nil)
		},
		func(err error) {
			require.NoError(t, err)
			require.FileExists(t, outputFilename)

			outputContent, err := os.ReadFile(outputFilename)
			require.NoError(t, err)
			assert.Equal(t, "ok", string(outputContent))
		},
	)
}

func TestClangTidy_Run(t *testing.T) {
	dataSet := [][]string{
		{},
		{"file1"},
	}

	for _, args := range dataSet {
		t.Run(fmt.Sprintf("args=%v", args), func(t *testing.T) {
			test.InvokeRunSuccess(t, clangTidyTool, internal.RunOptions{}, args, args)
		})
	}
}

func TestClangTidy_Run_Failed(t *testing.T) {
	test.InvokeRunFailed(t, clangTidyTool, internal.RunOptions{}, []string{}, []string{})
}

func clangTidyDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectClangTidy(ctx, options)
}

func clangTidyLinter(executable func() (*internal.Executable, error)) internal.ProjectLinter {
	return tool.NewClangTidy(executable)
}

func clangTidyTool(executable func() (*internal.Executable, error)) internal.Tool {
	return tool.NewClangTidy(executable)
}

func clangTidyProjectInfo(t *testing.T, filenames []string) internal.ProjectInfo {
	return internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: filenames,
					},
				},
			},
		},
	}
}
