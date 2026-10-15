package tool_test

import (
	"context"
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

func TestClangFormat_DetectClangFormat(t *testing.T) {
	clangFormat := test.RunDetectToolLastFound(t, clangFormatDetect, []string{"clang-format"})

	assert.Equal(t, tool.IdClangFormat, clangFormat.Id())
}

func TestClangFormat_DetectClangFormat_NotFound(t *testing.T) {
	clangFormat := test.RunDetectToolNotFound(t, clangFormatDetect, []string{"clang-format"})

	assert.Equal(t, tool.IdClangFormat, clangFormat.Id())
}

func TestClangFormat_DetectClangFormat_Config(t *testing.T) {
	clangFormat := test.RunDetectToolConfig(t, clangFormatDetect, "clang-format")

	assert.Equal(t, tool.IdClangFormat, clangFormat.Id())
}

func TestClangFormat_FormatFiles_All(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: internal.StrPtr("build"),
				StructureProvider: clangFormatProjectStructure(
					"target1",
					[]string{"file1.go", "file2.go"},
				),
			},
		},
		[]string{
			"--Werror",
			"-i",
			"file1.go",
			"file2.go",
		},
	)
}

func TestClangFormat_FormatFiles_All_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: internal.StrPtr("build"),
				StructureProvider: clangFormatProjectStructure(
					"target2",
					[]string{"file1.go", "file2.go"},
				),
			},
		},
		[]string{
			"--Werror",
			"-i",
			"file1.go",
			"file2.go",
		},
	)
}

func TestClangFormat_FormatFiles_All_CustomConfig(t *testing.T) {
	directory := t.TempDir()

	file, err := os.Create(filepath.Join(directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       directory,
				OutputDirectory: internal.StrPtr("build"),
				StructureProvider: clangFormatProjectStructure(
					"target3",
					[]string{"file1.go", "file2.go"},
				),
			},
		},
		[]string{
			fmt.Sprintf("--style=file:%v", filepath.Join(directory, ".clang-format")),
			"--Werror",
			"-i",
			"file1.go",
			"file2.go",
		},
	)
}

func TestClangFormat_FormatFiles(t *testing.T) {
	const directory = "project"

	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       directory,
				OutputDirectory: internal.StrPtr("build"),
			},
			Filenames: &[]string{"file1.go", filepath.Join(directory, "file3.go")},
		},
		[]string{
			"--Werror",
			"-i",
			"file1.go",
			filepath.Join(directory, "file3.go"),
		},
	)
}

func TestClangFormat_FormatFiles_Failed(t *testing.T) {
	const directory = "project"

	test.InvokeFormatFilesFailed(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       directory,
				OutputDirectory: internal.StrPtr("build"),
			},
			Filenames: &[]string{"file1.go", filepath.Join(directory, "file3.go")},
		},
		[]string{
			"--Werror",
			"-i",
			"file1.go",
			filepath.Join(directory, "file3.go"),
		},
	)
}

func TestClangFormat_FormatFiles_CustomConfig(t *testing.T) {
	directory := t.TempDir()

	file, err := os.Create(filepath.Join(directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       directory,
				OutputDirectory: internal.StrPtr("build"),
			},
			Filenames: &[]string{"file1.go", filepath.Join(directory, "file3.go")},
		},
		[]string{
			fmt.Sprintf("--style=file:%v", filepath.Join(directory, ".clang-format")),
			"--Werror",
			"-i",
			"file1.go",
			filepath.Join(directory, "file3.go"),
		},
	)
}

func TestClangFormat_FormatFiles_CheckAll(t *testing.T) {
	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: internal.StrPtr("build"),
				StructureProvider: clangFormatProjectStructure(
					"target4",
					[]string{"file1.go", "file2.go"},
				),
			},
			CheckOnly: true,
		},
		[]string{
			"--Werror",
			"--dry-run",
			"file1.go",
			"file2.go",
		},
	)
}

func TestClangFormat_FormatFiles_CheckAll_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       "project",
				OutputDirectory: internal.StrPtr("build"),
				StructureProvider: clangFormatProjectStructure(
					"target5",
					[]string{"file1.go", "file2.go"},
				),
			},
			CheckOnly: true,
		},
		[]string{
			"--Werror",
			"--dry-run",
			"file1.go",
			"file2.go",
		},
	)
}

func TestClangFormat_FormatFiles_CheckAll_CustomConfig(t *testing.T) {
	directory := t.TempDir()

	file, err := os.Create(filepath.Join(directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       directory,
				OutputDirectory: internal.StrPtr("build"),
				StructureProvider: clangFormatProjectStructure(
					"target1",
					[]string{"file1.go", "file2.go"},
				),
			},
			CheckOnly: true,
		},
		[]string{
			fmt.Sprintf("--style=file:%v", filepath.Join(directory, ".clang-format")),
			"--Werror",
			"--dry-run",
			"file1.go",
			"file2.go",
		},
	)
}

func TestClangFormat_FormatFiles_Check(t *testing.T) {
	const directory = "project"

	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       directory,
				OutputDirectory: internal.StrPtr("build"),
			},
			CheckOnly: true,
			Filenames: &[]string{"file1.go", filepath.Join(directory, "file3.go")},
		},
		[]string{
			"--Werror",
			"--dry-run",
			"file1.go",
			filepath.Join(directory, "file3.go"),
		},
	)
}

func TestClangFormat_FormatFiles_Check_Failed(t *testing.T) {
	test.InvokeFormatFilesFailed(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{Directory: "."},
			CheckOnly:   true,
			Filenames:   &[]string{"file1.go"},
		},
		[]string{
			"--Werror",
			"--dry-run",
			filepath.Join(".", "file1.go"),
		},
	)
}

func TestClangFormat_FormatFiles_Check_CustomConfig(t *testing.T) {
	directory := t.TempDir()

	file, err := os.Create(filepath.Join(directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	test.InvokeFormatFilesSuccess(
		t,
		clangFormatFormatter,
		internal.ProjectFormatterOptions{
			ProjectInfo: internal.ProjectInfo{
				Directory:       directory,
				OutputDirectory: internal.StrPtr("build"),
			},
			CheckOnly: true,
			Filenames: &[]string{"file1.go"},
		},
		[]string{
			fmt.Sprintf("--style=file:%v", filepath.Join(directory, ".clang-format")),
			"--Werror",
			"--dry-run",
			"file1.go",
		},
	)
}

func TestClangFormat_Run(t *testing.T) {
	dataSet := [][]string{
		{},
		{"file1"},
	}

	for _, args := range dataSet {
		t.Run(fmt.Sprintf("data=%v", args), func(t *testing.T) {
			test.InvokeRunSuccess(t, clangFormatTool, internal.RunOptions{}, args, args)
		})
	}
}

func TestClangFormat_Run_Failed(t *testing.T) {
	test.InvokeRunFailed(t, clangFormatTool, internal.RunOptions{}, []string{}, []string{})
}

func clangFormatDetect(ctx context.Context, options internal.DetectOptions) internal.Tool {
	return tool.DetectClangFormat(ctx, options)
}

func clangFormatFormatter(
	executable func() (*internal.Executable, error),
) internal.ProjectFormatter {
	return tool.NewClangFormat(executable)
}

func clangFormatTool(
	executable func() (*internal.Executable, error),
) internal.Tool {
	return tool.NewClangFormat(executable)
}

func clangFormatProjectStructure(
	target string,
	filenames []string,
) internal.ProjectStructureProvider {
	return &internal.FixedProjectStructureProvider{
		ProjectStructure: internal.ProjectStructure{
			Targets: map[string]internal.ProjectTarget{
				target: {
					Files: filenames,
				},
			},
		},
	}
}
