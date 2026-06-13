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

func TestClangFormat_DetectClangFormat(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(env.NewExecutable("/bin/clang-format"), nil)

	clangFormat := tool.DetectClangFormat(t.Context(), internal.DetectOptions{Environment: env})
	assert.NotNil(t, clangFormat)
	assert.Equal(t, "clang-format", clangFormat.Id())
	assert.True(t, clangFormat.IsAvailable())

	if executable := clangFormat.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/clang-format", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(nil, nil)

	clangFormat := tool.DetectClangFormat(t.Context(), internal.DetectOptions{Environment: env})
	assert.NotNil(t, clangFormat)
	assert.Equal(t, "clang-format", clangFormat.Id())
	assert.False(t, clangFormat.IsAvailable())
	assert.Nil(t, clangFormat.Executable())

	env.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_Config(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format-19").
		Return(env.NewExecutable("/bin/clang-format"), nil)

	clangFormat := tool.DetectClangFormat(t.Context(), internal.DetectOptions{
		Environment: env,
		ToolsPaths:  map[string]string{"clang-format": "clang-format-19"},
	})
	assert.NotNil(t, clangFormat)
	assert.Equal(t, "clang-format", clangFormat.Id())
	assert.True(t, clangFormat.IsAvailable())

	if executable := clangFormat.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/clang-format", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_All(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       "project",
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		"file1.go",
		"file2.go",
	}).
		Return(nil)

	err := clangFormat.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_All_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       "project",
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		"file1.go",
		"file2.go",
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_All_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	file, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"-i",
		"file1.go",
		"file2.go",
	}).
		Return(nil)

	err = clangFormat.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		"file1.go",
		filepath.Join(info.Directory, "file3.go"),
	}).
		Return(nil)

	err := clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{
			ProjectInfo: info,
			Filenames:   &[]string{"file1.go", filepath.Join(info.Directory, "file3.go")},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		"file1.go",
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, Filenames: &[]string{"file1.go"}},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
	}

	file, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"-i",
		"file1.go",
	}).
		Return(nil)

	err = clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, Filenames: &[]string{"file1.go"}},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       "project",
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		"file1.go",
		"file2.go",
	}).
		Return(nil)

	err := clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, CheckOnly: true},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CheckAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       "project",
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		"file1.go",
		"file2.go",
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, CheckOnly: true},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CheckAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
		StructureProvider: &internal.FixedProjectStructureProvider{
			ProjectStructure: internal.ProjectStructure{
				Targets: map[string]internal.ProjectTarget{
					"target1": {
						Files: []string{"file1.go", "file2.go"},
					},
				},
			},
		},
	}

	file, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"--dry-run",
		"file1.go",
		"file2.go",
	}).
		Return(nil)

	err = clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{ProjectInfo: info, CheckOnly: true},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Check(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		"file1.go",
		filepath.Join(info.Directory, "file3.go"),
	}).
		Return(nil)

	err := clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{
			ProjectInfo: info,
			CheckOnly:   true,
			Filenames:   &[]string{"file1.go", filepath.Join(info.Directory, "file3.go")},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Check_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{
			ProjectInfo: info,
			CheckOnly:   true,
			Filenames:   &[]string{"file1.go"},
		},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Check_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:       t.TempDir(),
		OutputDirectory: internal.StrPtr("build"),
	}

	file, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)
	assert.NoError(t, file.Close())

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"--dry-run",
		"file1.go",
	}).
		Return(nil)

	err = clangFormat.FormatFiles(
		t.Context(),
		internal.ProjectFormatterOptions{
			ProjectInfo: info,
			CheckOnly:   true,
			Filenames:   &[]string{"file1.go"},
		},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{}).
		Return(nil)

	err := clangFormat.RunForProject(t.Context(), desc, []string{})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{}).
		Return(errors.New("failed"))

	err := clangFormat.RunForProject(t.Context(), desc, []string{})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
