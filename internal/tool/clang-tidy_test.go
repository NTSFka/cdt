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

	clangTidy := tool.DetectClangTidy(t.Context(), env)
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

	clangTidy := tool.DetectClangTidy(t.Context(), env)
	assert.NotNil(t, clangTidy)
	assert.Equal(t, "clang-tidy", clangTidy.Id())
	assert.False(t, clangTidy.IsAvailable())
	assert.Nil(t, clangTidy.Executable())

	env.AssertExpectations(t)
}

func TestClangTidy_LintAll(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

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

	exec.OnRun("clang-tidy", []string{
		"-p", *info.IntermediateDirectory,
		"file1.go",
		"file2.go",
	}).
		Return(nil)

	err := clangTidy.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintAll_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

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

	exec.OnRun("clang-tidy", []string{
		"-p", *info.IntermediateDirectory,
		"file1.go",
		"file2.go",
	}).
		Return(errors.New("failed"))

	err := clangTidy.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintAll_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

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

	_, err := os.Create(filepath.Join(info.Directory, ".clang-tidy"))
	require.NoError(t, err)

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(info.Directory, ".clang-tidy")),
		"-p", *info.IntermediateDirectory,
		"file1.go",
		"file2.go",
	}).
		Return(nil)

	err = clangTidy.LintAll(t.Context(), internal.ProjectLinterOptions{ProjectInfo: info})
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:             t.TempDir(),
		IntermediateDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *info.IntermediateDirectory,
		"file1.go",
		filepath.Join(info.Directory, "file3.go"),
	}).
		Return(nil)

	err := clangTidy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"file1.go", filepath.Join(info.Directory, "file3.go")},
	)
	require.NoError(t, err)

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_Failed(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:             "project",
		IntermediateDirectory: internal.StrPtr("build"),
	}

	exec.OnRun("clang-tidy", []string{
		"-p", *info.IntermediateDirectory,
		"file1.go",
	}).
		Return(errors.New("failed"))

	err := clangTidy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"file1.go"},
	)
	require.EqualError(t, err, "failed")

	exec.AssertExpectations(t)
}

func TestClangTidy_LintFiles_CustomConfig(t *testing.T) {
	exec := test.NewExecutable(t)

	clangTidy := tool.NewClangTidy(exec.LazyExecutable("clang-tidy"))

	info := internal.ProjectInfo{
		Directory:             t.TempDir(),
		IntermediateDirectory: internal.StrPtr("build"),
	}

	_, err := os.Create(filepath.Join(info.Directory, ".clang-tidy"))
	require.NoError(t, err)

	exec.OnRun("clang-tidy", []string{
		fmt.Sprintf("--config-file=%v", filepath.Join(info.Directory, ".clang-tidy")),
		"-p", *info.IntermediateDirectory,
		"file1.go",
	}).
		Return(nil)

	err = clangTidy.LintFiles(
		t.Context(),
		internal.ProjectLinterOptions{ProjectInfo: info},
		[]string{"file1.go"},
	)
	require.NoError(t, err)

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
