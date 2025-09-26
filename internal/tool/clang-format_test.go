package tool_test

import (
	"cdt/internal"
	"cdt/internal/test"
	"cdt/internal/tool"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClangFormat_DetectClangFormat(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(env.NewExecutable("/bin/clang-format"))

	clangFormat := tool.DetectClangFormat(t.Context(), env)
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
		Return(nil)

	clangFormat := tool.DetectClangFormat(t.Context(), env)
	assert.NotNil(t, clangFormat)
	assert.Equal(t, "clang-format", clangFormat.Id())
	assert.False(t, clangFormat.IsAvailable())
	assert.Nil(t, clangFormat.Executable())

	env.AssertExpectations(t)
}

func TestClangFormat_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
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
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(nil)

	err := clangFormat.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
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
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:             t.TempDir(),
		IntermediateDirectory: internal.StrPtr("build"),
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

	_, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(nil)

	err = clangFormat.FormatAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file3.go"),
	}).
		Return(nil)

	err := clangFormat.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go", filepath.Join(info.Directory, "file3.go")})
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
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	_, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(nil)

	err = clangFormat.FormatFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
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
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(nil)

	err := clangFormat.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
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
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{
		Directory:             t.TempDir(),
		IntermediateDirectory: internal.StrPtr("build"),
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

	_, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(nil)

	err = clangFormat.FormatCheckAll(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file3.go"),
	}).
		Return(nil)

	err := clangFormat.FormatCheckFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go", filepath.Join(info.Directory, "file3.go")})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(errors.New("failed"))

	err := clangFormat.FormatCheckFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangFormat := tool.NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	_, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	require.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(nil)

	err = clangFormat.FormatCheckFiles(t.Context(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
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
