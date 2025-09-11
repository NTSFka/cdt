package tool

import (
	"cdt/internal"
	"cdt/internal/test"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClangFormat_DetectClangFormat(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(env.NewExecutable("/bin/clang-format"))

	tool := DetectClangFormat(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.True(t, tool.IsAvailable())

	if executable := tool.Executable(); assert.NotNil(t, executable) {
		assert.Equal(t, "/bin/clang-format", executable.Path)
	}

	env.AssertExpectations(t)
}

func TestClangFormat_DetectClangFormat_NotFound(t *testing.T) {
	env := test.NewEnvironment(t)
	env.OnFindExecutable("clang-format").
		Return(nil)

	tool := DetectClangFormat(context.Background(), env)
	assert.NotNil(t, tool)
	assert.Equal(t, "clang-format", tool.Id())
	assert.False(t, tool.IsAvailable())
	assert.Nil(t, tool.Executable())

	env.AssertExpectations(t)
}

func TestClangFormat_FormatAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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

	err := tool.FormatAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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

	err := tool.FormatAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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
	assert.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(nil)

	err = tool.FormatAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file3.go"),
	}).
		Return(nil)

	err := tool.FormatFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go", filepath.Join(info.Directory, "file3.go")})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"-i",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	_, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	assert.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"-i",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(nil)

	err = tool.FormatFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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

	err := tool.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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

	err := tool.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

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
	assert.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file2.go"),
	}).
		Return(nil)

	err = tool.FormatCheckAll(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
		filepath.Join(info.Directory, "file3.go"),
	}).
		Return(nil)

	err := tool.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go", filepath.Join(info.Directory, "file3.go")})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(errors.New("failed"))

	err := tool.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangFormat_FormatCheckFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	info := internal.ProjectInfo{Directory: t.TempDir(), IntermediateDirectory: internal.StrPtr("build")}

	_, err := os.Create(filepath.Join(info.Directory, ".clang-format"))
	assert.NoError(t, err)

	exec.OnRun("clang-format", []string{
		fmt.Sprintf("--style=file:%v", filepath.Join(info.Directory, ".clang-format")),
		"--Werror",
		"--dry-run",
		filepath.Join(info.Directory, "file1.go"),
	}).
		Return(nil)

	err = tool.FormatCheckFiles(context.Background(), internal.ProjectFormatterOptions{ProjectInfo: info}, []string{"file1.go"})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_Run(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{}).
		Return(nil)

	err := tool.RunForProject(context.Background(), desc, []string{})
	assert.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangFormat_Run_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	tool := NewClangFormat(exec.LazyExecutable("clang-format"))

	desc := internal.ProjectInfo{Directory: "."}

	exec.OnRun("clang-format", []string{}).
		Return(errors.New("failed"))

	err := tool.RunForProject(context.Background(), desc, []string{})
	assert.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}
